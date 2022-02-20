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
	if v, ok := raw["target"]; !ok || v == nil {
		return fmt.Errorf("field target: required")
	}
	type Plain ExecRequestsElemDataPersistenceDataOutElem
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = ExecRequestsElemDataPersistenceDataOutElem(plain)
	return nil
}

type ExecRequestsElem struct {
	// The wrapper to define sticky data dependency and generation properties.
	DataPersistence *ExecRequestsElemDataPersistence `json:"data-persistence,omitempty"`

	// Delayed execution wait time before start - in seconds.
	DelaySeconds int `json:"delay-seconds"`

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
	*j = ExecRequestsElem(plain)
	return nil
}

// Perforamnce test execution configuration schema
type Exec struct {
	// Test duration in seconds
	DurationSeconds int `json:"duration-seconds"`

	// Key/value pairs, defined for the test runtime.
	Environment []*ExecEnvironmentElem `json:"environment,omitempty"`

	// Optional settings for influxdb reporting.
	InfluxdbSettings *ExecInfluxdbSettings `json:"influxdb-settings,omitempty"`

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
