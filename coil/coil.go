// The coil is the orchestrator and executor of the timeline's tasks.
// The timeline is the plan and the tasks are the points in time, which triggers the execution.
// Termination: The coild should stop only on statistical condition, never on a single failure.
package coil

import (
	"context"
	"log"
	"time"

	datastore "github.com/csabakissmalta/tpee/datastore"
	execconf "github.com/csabakissmalta/tpee/exec"
	natsconnect "github.com/csabakissmalta/tpee/nats_connect"
	request "github.com/csabakissmalta/tpee/request"
	sessionstore "github.com/csabakissmalta/tpee/sessionstore"
	task "github.com/csabakissmalta/tpee/task"
	timeline "github.com/csabakissmalta/tpee/timeline"
)

const (
	// DEFAULT: compares two timestamps next to each other and uses the duration between them to time the execution
	COMPARE_TIMESTAMPS_MODE = "compare-timestamps-mode"

	// runs a timer and executes, when the message is dispatched from the timer
	GO_TIMER_MODE = "go-timer-mode"
)

// sessionstore capacity - default value
var SESSION_STORE_CAPACITY int = sessionstore.STORE_CAPACITY

type Coil struct {
	NATSClient              *natsconnect.NATSClient
	Ctx                     context.Context
	Timelines               []*timeline.Timeline
	EnvVars                 []*execconf.ExecEnvironmentElem
	DataStore               *datastore.DataBroadcaster
	SessionStore            *sessionstore.Store
	ResultsReportingChannel chan *task.Task
	ExecutionMode           string
	ConsumeClock            *time.Ticker
	StartTime               time.Time
}

type Option func(*Coil)

func WithContext(ctx context.Context) Option {
	return func(c *Coil) {
		c.Ctx = ctx
	}
}

func WithTimelines(tls []*timeline.Timeline) Option {
	return func(c *Coil) {
		c.Timelines = tls
	}
}

func WithEnvVariables(ev []*execconf.ExecEnvironmentElem) Option {
	return func(c *Coil) {
		c.EnvVars = ev
	}
}

func WithResultsReportingChannel(ch chan *task.Task) Option {
	return func(c *Coil) {
		c.ResultsReportingChannel = ch
	}
}

func WithExecutionMode(em string) Option {
	return func(c *Coil) {
		c.ExecutionMode = em
	}
}

func WithNATSClient(nc *natsconnect.NATSClient) Option {
	return func(c *Coil) {
		c.NATSClient = nc
	}
}

func New(option ...Option) *Coil {
	c := &Coil{
		ExecutionMode: COMPARE_TIMESTAMPS_MODE,
	}
	for _, o := range option {
		o(c)
	}
	return c
}

// It is a consumer of tasks relying on the main loop and context.
// It controls only the exact execution of the timeline
// Should start always from the first element and progressively consume the tasks.
func (c *Coil) Start() {
	// Create datastore
	c.createDatastore()

	// Create sessionstore
	c.createSessionstore()

	for _, tLine := range c.Timelines {
		// Validate timeline data before start
		c.validateTimelineData(tLine)

		if c.ExecutionMode == COMPARE_TIMESTAMPS_MODE {
			// Start consuming compare mode
			c.consumeTimelineCompareMode(tLine, c.EnvVars, c.ResultsReportingChannel)
		} else if c.ExecutionMode == GO_TIMER_MODE {
			// Start consuming timer mode
			c.consumeTimelineTimerMode(tLine, c.EnvVars, c.ResultsReportingChannel)
		}
	}
	<-make(chan bool)
}

// Stops the coil loop
// conditions can be different:
// 1. out of tasks
// 2. exeception
// 3. condition based termination
func (c *Coil) Stop() error {
	return nil
}

// Datastore
func (c *Coil) createDatastore() {
	// create the datastore
	all_req_conf := []*execconf.ExecRequestsElem{}
	for _, t := range c.Timelines {
		all_req_conf = append(all_req_conf, t.Rules)
	}
	names := execconf.GetAllDataPersistenceDataNames(all_req_conf)
	if len(names) > 0 {
		c.DataStore = datastore.New(
			datastore.WithDataOutSocketNames(names),
		)
		go c.DataStore.StartConsumingDataIn()
	}
}

func (c *Coil) createSessionstore() {
	c.SessionStore = sessionstore.NewStore(
		sessionstore.WithInOutCapacity(SESSION_STORE_CAPACITY),
	)
	go c.SessionStore.Start()
}

// Validate timeline data
func (c *Coil) validateTimelineData(tl *timeline.Timeline) {
	// validate
	e := timeline.CheckPostmanRequestAndValidateRequirements(tl.RequestBlueprint, c.EnvVars)
	if e != nil {
		log.Fatalf("DATA ERROR: %s", e.Error())
	}
}

// The function is the engine's main task - run the test
func (c *Coil) consumeTimelineCompareMode(tl *timeline.Timeline, env []*execconf.ExecEnvironmentElem, res_ch chan *task.Task) {
	// The Coil needs to control timelines in a separate routines
	go func() {
		var after_rmpup_ts int = 0
		if tl.CurrectTask.PlannedExecTimeNanos > 0 {
			time.Sleep(time.Duration(tl.CurrectTask.PlannedExecTimeNanos * int(time.Nanosecond)))
		}

		select {
		case tl.CurrectTask = <-tl.RampupTasks:
		case tl.CurrectTask = <-tl.Tasks:
		}

		// compose/execute task here
		_, ses, _ := request.ComposeHttpRequest(tl.CurrectTask, *tl.RequestBlueprint, tl.Rules.DataPersistence.DataIn, tl.Rules, tl.Feeds, c.DataStore, c.SessionStore)
		tl.CurrectTask.Execute(tl.HTTPClient, tl.Rules.DataPersistence.DataOut, tl.Rules.DataPersistence.DataIn, env, res_ch, *tl.Rules.CreatesSession, c.SessionStore, c.DataStore, ses)

		var next *task.Task

		for {
			select {
			case next = <-tl.RampupTasks:
				after_rmpup_ts = tl.CurrectTask.PlannedExecTimeNanos
			default:
				next = <-tl.Tasks
			}
			dorm_period := (next.PlannedExecTimeNanos - tl.CurrectTask.PlannedExecTimeNanos + after_rmpup_ts) * int(time.Nanosecond)
			time.Sleep(time.Duration(dorm_period))

			// compose/execute task here
			_, ses, _ := request.ComposeHttpRequest(next, *tl.RequestBlueprint, tl.Rules.DataPersistence.DataIn, tl.Rules, tl.Feeds, c.DataStore, c.SessionStore)
			next.Execute(tl.HTTPClient, tl.Rules.DataPersistence.DataOut, tl.Rules.DataPersistence.DataIn, env, res_ch, *tl.Rules.CreatesSession, c.SessionStore, c.DataStore, ses)
			tl.CurrectTask = next
		}
	}()
}

func (c *Coil) consumeTimelineTimerMode(tl *timeline.Timeline, env []*execconf.ExecEnvironmentElem, res_ch chan *task.Task) {
	// Set the timer with the duration of the step size
	c.ConsumeClock = time.NewTicker(time.Duration(tl.StepDuration * int(time.Nanosecond)))
	done := make(chan bool)

	if tl.Rules.DelaySeconds > 0 {
		time.Sleep(time.Duration(tl.Rules.DelaySeconds * int(time.Second)))
	}

	var next *task.Task
	var corr int64
	c.StartTime = time.Now()

	// Start the timer
	go func() {
		if len(tl.RampupTasks) > 0 {
			tl.CurrectTask = <-tl.RampupTasks
			// compose/execute task here
			// _, ses, _ := request.ComposeHttpRequest(tl.CurrectTask, *tl.RequestBlueprint, tl.Rules.DataPersistence.DataIn, tl.Rules, tl.Feeds, c.DataStore, c.SessionStore)
			// tl.CurrectTask.Execute(tl.HTTPClient, tl.Rules.DataPersistence.DataOut, tl.Rules.DataPersistence.DataIn, env, res_ch, *tl.Rules.CreatesSession, c.SessionStore, c.DataStore, ses)
		}

		for {
			select {
			case next = <-tl.RampupTasks:
				// if there is rampup, it falls back to compare mode
				elapsedTotal := time.Since(c.StartTime).Nanoseconds()
				corr = elapsedTotal - int64(tl.CurrectTask.PlannedExecTimeNanos)
				planned_delta := next.PlannedExecTimeNanos - tl.CurrectTask.PlannedExecTimeNanos
				dorm_period := planned_delta - int(corr)
				if dorm_period < 0 {
					dorm_period = 0
				}

				time.Sleep(time.Duration(dorm_period))

				// compose/execute task here
				_, ses, _ := request.ComposeHttpRequest(next, *tl.RequestBlueprint, tl.Rules.DataPersistence.DataIn, tl.Rules, tl.Feeds, c.DataStore, c.SessionStore)
				next.Execute(tl.HTTPClient, tl.Rules.DataPersistence.DataOut, tl.Rules.DataPersistence.DataIn, env, res_ch, *tl.Rules.CreatesSession, c.SessionStore, c.DataStore, ses)
				tl.CurrectTask = next

				// rampupStopwatch = time.Since(testStartTime)
			default:
				select {
				case <-done:
					return
				case <-c.ConsumeClock.C:
					// compose/execute task here
					next = <-tl.Tasks
					_, ses, _ := request.ComposeHttpRequest(next, *tl.RequestBlueprint, tl.Rules.DataPersistence.DataIn, tl.Rules, tl.Feeds, c.DataStore, c.SessionStore)
					next.Execute(tl.HTTPClient, tl.Rules.DataPersistence.DataOut, tl.Rules.DataPersistence.DataIn, env, res_ch, *tl.Rules.CreatesSession, c.SessionStore, c.DataStore, ses)
					tl.CurrectTask = next
				}
			}
		}
	}()
}

// --- NATS implementation for modifying traffic rate ---
func (c *Coil) UpdateTrafficRateFromNATSKVUpdate(tl_name string, tr *timeline.Transition, orig_tl_dur int) error {
	// impl
	go func() {
		var tln *timeline.Timeline

		// filter the timeline we need to change the rate
		for _, tl := range c.Timelines {
			if tl.Name == tl_name {
				tln = tl
				break
			}
		}

		tln.Repopulate(tr, orig_tl_dur, c.StartTime, c.ConsumeClock)
	}()

	return nil
}
