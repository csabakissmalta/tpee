# p-tpee
Performance test plan execution engine

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