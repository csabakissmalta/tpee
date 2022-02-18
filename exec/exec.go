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

func GetAllDataPersistenceDataNames(lst []*ExecRequestsElem) []string {
	names := []string{}
	c := 0
	for _, le := range lst {
		for _, l := range le.DataPersistence.DataOut {
			names[c] = *l.Name
			c += 1
		}
	}
	return names
}
