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
	request "github.com/csabakissmalta/tpee/request"
	task "github.com/csabakissmalta/tpee/task"
	timeline "github.com/csabakissmalta/tpee/timeline"
)

const (
	// DEFAULT: compares two timestamps next to each other and uses the duration between them to time the execution
	COMPARE_TIMESTAMPS_MODE = "compare-timestamps-mode"

	// runs a timer and executes, when the message is dispatched from the timer
	GO_TIMER_MODE = "go-timer-mode"
)

type Coil struct {
	Ctx                     context.Context
	Timelines               []*timeline.Timeline
	EnvVars                 []*execconf.ExecEnvironmentElem
	DataStore               *datastore.DataBroadcaster
	ResultsReportingChannel chan *task.Task
	ExecutionMode           string
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
	if c.ExecutionMode == COMPARE_TIMESTAMPS_MODE {
		for _, tLine := range c.Timelines {
			c.consumeTimelineCompareMode(tLine, c.EnvVars, c.ResultsReportingChannel)
		}
	} else if c.ExecutionMode == GO_TIMER_MODE {
		for _, tLine := range c.Timelines {
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

// The function is the engine's main task - run the test
func (c *Coil) consumeTimelineCompareMode(tl *timeline.Timeline, env []*execconf.ExecEnvironmentElem, res_ch chan *task.Task) {
	// create the datastore
	if c.DataStore == nil {
		all_req_conf := []*execconf.ExecRequestsElem{}
		for _, t := range c.Timelines {
			all_req_conf = append(all_req_conf, t.Rules)
		}
		names := execconf.GetAllDataPersistenceDataNames(all_req_conf)
		c.DataStore = datastore.New(
			datastore.WithDataOutSocketNames(names),
		)
		go c.DataStore.StartConsumingDataIn()
	}

	e := timeline.CheckPostmanRequestAndValidateRequirements(tl.RequestBlueprint, c.EnvVars)
	if e != nil {
		log.Fatalf("DATA ERROR: %s", e.Error())
	}

	// The Coil needs to control timelines in a separate routines
	go func() {
		tl.CurrectTask = <-tl.Tasks
		if tl.CurrectTask.PlannedExecTimeNanos > 0 {
			time.Sleep(time.Duration(tl.CurrectTask.PlannedExecTimeNanos * int(time.Nanosecond)))
		}
		// compose/execute task here
		request.ComposeHttpRequest(tl.CurrectTask, *tl.RequestBlueprint, env, tl.Feeds, c.DataStore)
		tl.CurrectTask.Execute(tl.HTTPClient, tl.Rules.DataPersistence.DataOut, res_ch)

		for {
			next := <-tl.Tasks
			dorm_period := (next.PlannedExecTimeNanos - tl.CurrectTask.PlannedExecTimeNanos) * int(time.Nanosecond)
			time.Sleep(time.Duration(dorm_period))

			// compose/execute task here
			request.ComposeHttpRequest(next, *tl.RequestBlueprint, env, tl.Feeds, c.DataStore)
			next.Execute(tl.HTTPClient, tl.Rules.DataPersistence.DataOut, res_ch)
			tl.CurrectTask = next
		}
	}()
}

func (c *Coil) consumeTimelineTimerMode(tl *timeline.Timeline, env []*execconf.ExecEnvironmentElem, res_ch chan *task.Task) {
	// Set the timer with the duration of the step size
	engine_timer := *time.NewTimer(time.Duration(tl.StepDuration * int(time.Nanosecond)))

	// Start the timer
	go func() {
		c := 0
		for {
			<-engine_timer.C
			log.Println("EXEC TASK")
			c++
			if c == len(tl.Tasks) {
				break
			}
		}
	}()
}
