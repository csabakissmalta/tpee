package postman

import "io/ioutil"

func (p *Postman) LoadCollection(path string) error {
	fb, e := ioutil.ReadFile(path)
	if e != nil {
		return e
	}
	e = p.UnmarshalJSON(fb)
	if e != nil {
		return e
	}
	return nil
}
