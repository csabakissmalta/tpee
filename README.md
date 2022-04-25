# tpee (Performance test plan execution engine)
Collection of modules

## Basic sub-modules and functionality
This is a collection of sub-modules, which work together to generate traffic, using Golang's http client and Postman collection as a definition base. 

### coil
The main controller. Has two running modes:
1. `COMPARE_TIMESTAMPS_MODE` - is the mode, when task execution relies on the `Task`'s pre-defined timestamp, which is a relative time in nanoseconds, from the start of the test.
2. `GO_TIMER_MODE` - Tasks are executed in a monotonic way, based on a `Ticker`'s events.

It can control multiple timelines (plans). The duration is controlled via context. When the context expires, the `coil` stops the execution.

### timeline
It is the actual blueprint of a plan, how a request is going to be executed during the test. The timeline has the definition of the rate and validates and controls the request definition data. It also takes care of the feeds, which are loaded into data structures and fed into the requests runtime.

### task
These are the atomic elements of the test. As it is described above, the task has the definition of the planned execution time, and when it is executed, the absolute time of execution. Contain request-related data and references to the http request and response objects. These objects are the holders of all data stored or reported. 

### request
Is a sub-structure of the task, which is built or finilised runtime, based on data, required to passed into it. 

## Extended sub-modules and their role
The following sub-modules are defining the test properties and execution rules.

### exec
This is the loader and parser of the execution configuration, which is a separate json file. The json is defined by a schema, which is stored in the repository for the module's implementation. The structures and parsers are also generaed from that schema and any change in the schema requires re-generation of the `exec-schema.go` file as well.

### postman
As the name indicates, it contains the Postman collection-related code, which is similarly as the `exec` package has a loading and parsing functionality and also generated from Postman's json schema, which is publicly available. 

### datastore
Is a datastore, which can contain data, not bound to sessions. It's basic structure contains a slice of channels, identified by a name and can be consumed by tasks/requests.

### sessionstore
A ring buffer, with a pre-defined capacity, which also controls the sessions time-based validity, when it is being consumed. Possible to attach data to it. Stores client session cookies, which are then set on the requests runtime, when consumed from the store, and fed vback to store, after usage.

## Example usage
``` go
import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"time"

	engine "github.com/csabakissmalta/tpee/coil"
	execconf "github.com/csabakissmalta/tpee/exec"
	postman "github.com/csabakissmalta/tpee/postman"
	task "github.com/csabakissmalta/tpee/task"
	timeline "github.com/csabakissmalta/tpee/timeline"
)

// Engine instance
var runner *engine.Coil

// Collection
var postman_coll *postman.Postman

// Timelines
var timelines []*timeline.Timeline

// Exec Config
var exec_conf *execconf.Exec

// Reporting channel
var rep_chan chan *task.Task

func init() {
	// Check flags

	// Load execution config
	exec_conf = &execconf.Exec{}
	exec_conf.LoadExecConfig("conf/exec_config.json")

	// Load postman collection
	postman_coll = &postman.Postman{}
	postman_coll.LoadCollection("postman/postman_collection.json")

	// Create timelines based on the config
	timelines = make([]*timeline.Timeline, len(exec_conf.Requests))
	for i := 0; i < len(exec_conf.Requests); i++ {
		req_rules := exec_conf.Requests[i]

		t := http.DefaultTransport.(*http.Transport).Clone()
		t.MaxIdleConns = 100
		t.MaxConnsPerHost = 100
		t.MaxIdleConnsPerHost = 100
		t.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

		client := http.Client{
			Transport: t,
			Timeout:   time.Duration(3 * time.Second),
		}
		tmln := timeline.New(
			timeline.WithRules(req_rules),
			timeline.WithHTTPClient(&client),
		)
		postmanReq, e := postman_coll.GetRequestByName(req_rules.Name)
		if e != nil {
			continue
		}
		tmln.Populate(exec_conf.DurationSeconds, postmanReq, exec_conf.Environment)
		timelines[i] = tmln
	}

	// Reporting channel
	rep_chan = make(chan *task.Task, 1000)

	// Create engine
	runner = engine.New(
		engine.WithContext(context.Background()),
		engine.WithTimelines(timelines),
		engine.WithEnvVariables(exec_conf.Environment),
		engine.WithResultsReportingChannel(rep_chan),
	)
}

func main() {
	if rep_chan != nil {
		go func() {
			for {
				t := <-rep_chan
				log.Println(t.Request.URL, "\n", t.Response.StatusCode)
			}
		}()
	}
	runner.Start()
}
```