package exec

import "io/ioutil"

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
