package postman

import (
	"errors"
	"io/ioutil"
)

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

func (p *Postman) GetRequestByName(name string) (*Request, error) {
	for _, r := range p.Items {
		if name == *r.Name {
			return &r.Request, nil
		}
	}
	return nil, errors.New("no request in the postman collection with that name")
}
