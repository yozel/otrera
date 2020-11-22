package template

import (
	"text/template"
)

var funcMap map[string]interface{} = template.FuncMap{
	"replace": templateFuncReplace,
	"gjson":   templateFuncGjson,
}

// New creates a new template with given templateString
func New(templateString string) (*template.Template, error) {
	return template.New("").Funcs(funcMap).Parse(templateString)
}
