package exec

import (
	"encoding/json"
	"fmt"
)

type ExecEnvironmentElem struct {
	// The env variable name
	Key string `json:"key"`

	// Enum for setting variables for later parsing and composition of the executable
	// http.Request
	Type *string `json:"type,omitempty"`

	// The env variable value
	Value string `json:"value"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ExecEnvironmentElem) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["key"]; !ok || v == nil {
		return fmt.Errorf("field key: required")
	}
	if v, ok := raw["value"]; !ok || v == nil {
		return fmt.Errorf("field value: required")
	}
	type Plain ExecEnvironmentElem
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ExecEnvironmentElem(plain)
	return nil
}

// Optional settings for influxdb reporting.
type ExecInfluxdbSettings struct {
	// The database/bucket name.
	Database string `json:"database"`

	// The database host.
	DatabaseHost string `json:"database-host"`

	// The organisation set for the database instance.
	DatabaseOrg string `json:"database-org"`

	// Influxdb2 uses tokens, this is the one to connect to the db.
	DatabaseToken string `json:"database-token"`

	// Measurement for the test data.
	Measurement string `json:"measurement"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ExecInfluxdbSettings) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["database"]; !ok || v == nil {
		return fmt.Errorf("field database: required")
	}
	if v, ok := raw["database-host"]; !ok || v == nil {
		return fmt.Errorf("field database-host: required")
	}
	if v, ok := raw["database-org"]; !ok || v == nil {
		return fmt.Errorf("field database-org: required")
	}
	if v, ok := raw["database-token"]; !ok || v == nil {
		return fmt.Errorf("field database-token: required")
	}
	if v, ok := raw["measurement"]; !ok || v == nil {
		return fmt.Errorf("field measurement: required")
	}
	type Plain ExecInfluxdbSettings
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ExecInfluxdbSettings(plain)
	return nil
}

// The precise definition for the implementation where to find the required value
// to save.
type ExecRequestsElemDataPersistenceDataOutElem struct {
	// Should be set to be able to determine the way to extract the value.
	ContentType string `json:"content-type"`

	// The name of the property the value of is required
	Name *string `json:"name,omitempty"`

	// Target corresponds to the JSON schema field "target".
	Target interface{} `json:"target"`

	// Storage corresponds to the JSON schema field "storage".
	Storage interface{} `json:"storage,omitempty"`

	// Whether the objects need to be disposed after usage.
	Retention bool `json:"retention,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ExecRequestsElemDataPersistenceDataOutElem) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["content-type"]; !ok || v == nil {
		return fmt.Errorf("field content-type: required")
	}
	if v, ok := raw["name"]; !ok || v == nil {
		return fmt.Errorf("field name: required")
	}
	if v, ok := raw["storage"]; !ok || v == nil {
		return fmt.Errorf("field storage: required")
	}
	if v, ok := raw["target"]; !ok || v == nil {
		return fmt.Errorf("field target: required")
	}
	type Plain ExecRequestsElemDataPersistenceDataOutElem
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if v, ok := raw["retention"]; !ok || v == nil {
		plain.Retention = false
	}
	*j = ExecRequestsElemDataPersistenceDataOutElem(plain)
	return nil
}

type ExecRequestsElem struct {
	// Definition for the request, whether creates or uses session
	CreatesSession *bool `json:"creates-session,omitempty"`

	// The wrapper to define sticky data dependency and generation properties.
	DataPersistence *ExecRequestsElemDataPersistence `json:"data-persistence,omitempty"`

	// Delayed execution wait time before start - in seconds.
	DelaySeconds int `json:"delay-seconds"`

	// Set the client to follow the redirects or not.
	FollowRedirects bool `json:"follow-redirects,omitempty"`

	// Per second execution rate.
	Frequency int `json:"frequency"`

	// The request's name from the Postman collection.
	Name string `json:"name"`

	// Order number of the request. Can be set for maintaining data dependency.
	OrderNumber int `json:"order-number"`
}

// The wrapper to define sticky data dependency and generation properties.
type ExecRequestsElemDataPersistence struct {
	// Data variable names, what the request is dependant on.
	DataIn []string `json:"data-in,omitempty"`

	// Data variable names, generated/set from the request/response
	DataOut []*ExecRequestsElemDataPersistenceDataOutElem `json:"data-out,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ExecRequestsElem) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["delay-seconds"]; !ok || v == nil {
		return fmt.Errorf("field delay-seconds: required")
	}
	if v, ok := raw["frequency"]; !ok || v == nil {
		return fmt.Errorf("field frequency: required")
	}
	if v, ok := raw["name"]; !ok || v == nil {
		return fmt.Errorf("field name: required")
	}
	if v, ok := raw["order-number"]; !ok || v == nil {
		return fmt.Errorf("field order-number: required")
	}
	type Plain ExecRequestsElem
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if v, ok := raw["follow-redirects"]; !ok || v == nil {
		plain.FollowRedirects = true
	}
	*j = ExecRequestsElem(plain)
	return nil
}

// The HDR Histogram output settings.
type ExecHdrHistogramSettings struct {
	// The base path, where the files should be saved
	BaseOutPath string `json:"base-out-path"`

	// The additional identifier for the set of files from the test. Can be the tesyed
	// version of subject.
	VersionLabel string `json:"version-label"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ExecHdrHistogramSettings) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["base-out-path"]; !ok || v == nil {
		return fmt.Errorf("field base-out-path: required")
	}
	if v, ok := raw["version-label"]; !ok || v == nil {
		return fmt.Errorf("field version-label: required")
	}
	type Plain ExecHdrHistogramSettings
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ExecHdrHistogramSettings(plain)
	return nil
}

// Period, which starts at 0 and reaches the timelines traffic level.
type ExecRampup struct {
	// Rampup period duration.
	DurationSeconds *int `json:"duration-seconds,omitempty"`

	// RampupType corresponds to the JSON schema field "rampup-type".
	RampupType *string `json:"rampup-type,omitempty"`
}

// Perforamnce test execution configuration schema
type Exec struct {
	// Test duration in seconds
	DurationSeconds int `json:"duration-seconds"`

	// Key/value pairs, defined for the test runtime.
	Environment []*ExecEnvironmentElem `json:"environment,omitempty"`

	// Optional settings for influxdb reporting.
	InfluxdbSettings *ExecInfluxdbSettings `json:"influxdb-settings,omitempty"`

	// Period, which starts at 0 and reaches the timelines traffic level.
	Rampup *ExecRampup `json:"rampup,omitempty"`

	// The HDR Histogram output settings.
	HdrHistogramSettings *ExecHdrHistogramSettings `json:"hdr-histogram-settings,omitempty"`

	// The requests and their rate definition corresponding with the Postman
	// collection
	Requests []*ExecRequestsElem `json:"requests"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Exec) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["duration-seconds"]; !ok || v == nil {
		return fmt.Errorf("field duration-seconds: required")
	}
	if v, ok := raw["requests"]; !ok || v == nil {
		return fmt.Errorf("field requests: required")
	}
	type Plain Exec
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = Exec(plain)
	return nil
}
