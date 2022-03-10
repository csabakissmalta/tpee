package exec

import "io/ioutil"

const (
	STRING_VALUE    = "stringValue"
	FEED_VALUE      = "feedValue"
	GENERATED_VALUE = "generatedValue"
	SESSION_VALUE   = "sessionValue"
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
	for _, le := range lst {
		if le.DataPersistence != nil && len(le.DataPersistence.DataOut) > 0 {
			for _, l := range le.DataPersistence.DataOut {
				names = append(names, *l.Name)
			}
		}
	}
	return names
}
