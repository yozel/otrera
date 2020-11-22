package template

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/tidwall/gjson"
)

func templateFuncReplace(find, replace, input string) string {
	return strings.Replace(input, find, replace, -1)
}

func templateFuncGjson(query string, input []interface{}) interface{} {
	j, err := json.Marshal(input)
	if err != nil {
		log.Printf("%s\n", err)
		return ""
	}
	return gjson.Get(string(j), query)
}
