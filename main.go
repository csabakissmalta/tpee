package main

import (
	"io/ioutil"
	"log"

	"github.com/csabakissmalta/tpee/request"
)

func load_coll(p string) request.PostmanSchemaJson {
	bts, e := ioutil.ReadFile(p)
	if e != nil {
		log.Fatal("couldn't load the file")
	}

	collection := request.PostmanSchemaJson{}
	collection.UnmarshalJSON(bts)
	return collection
}

func main() {
	c := load_coll("postman_collections/crs_collection.json")
	for _, i := range c.Item {
		// itype := reflect.TypeOf(i)
		// log.Println(itype)
		// numFields := itype.NumField()
		// ri := reflect.ValueOf(&i)
		// for j := 0; j < numFields; j++ {
		// 	log.Println(ri.Elem().Field(j))
		// }
		for k := range i.(map[string]interface{}) {
			if k == "request" {
				r := i.(map[string]interface{})[k].(request.Request)
				log.Println(r)
			}
		}
		log.Println("---")
	}
}
