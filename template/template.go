package template

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
	"text/template"

	"github.com/tidwall/gjson"
	"github.com/yozel/otrera/objectstore"
)

var funcMap map[string]interface{} = template.FuncMap{
	"replace": func(find, replace, input string) string {
		return strings.Replace(input, find, replace, -1)
	},
	"gjson": func(query string, input interface{}) interface{} {
		j, err := json.Marshal(input)
		if err != nil {
			log.Printf("%s\n", err)
			return ""
		}
		r := gjson.Get(string(j), query)
		if r.Exists() {
			return r
		} else {
			return ""
		}

	},
	"typeof": reflect.TypeOf,
}

// New creates a new template with given templateString
func New(name, templateString string, s *objectstore.ObjectStore) (*template.Template, error) {
	var dynamicFuncMap map[string]interface{} = template.FuncMap{
		"getall": func(key string, labels ...string) ([]objectstore.Object, error) {
			r, err := s.GetAll(key, nil) // TODO: support labels
			if err != nil {
				return nil, err
			}
			return r, nil
		},
		"get": func(key string, keyAppend ...string) (*objectstore.Object, error) {
			r, err := s.Get(fmt.Sprintf("%s%s", key, strings.Join(keyAppend, "")))
			if err != nil {
				return nil, err
			}
			return r, nil
		},
	}

	return template.New(name).Funcs(funcMap).Funcs(dynamicFuncMap).Parse(templateString)
}
