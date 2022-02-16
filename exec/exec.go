package exec

import "io/ioutil"

const (
	STRING_VALUE    = "stringValue"
	FEED_VALUE      = "feedValue"
	GENERATED_VALUE = "generatedValue"
)

func (ex *Exec) LoadExecConfig(path string) error {
	fb, e := ioutil.ReadFile(path)
	if e != nil {
		return e
	}
	e = ex.UnmarshalJSON(fb)
	if e != nil {
		return e
	}
	return nil
}
