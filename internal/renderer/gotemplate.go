package renderer

import (
	"bytes"

	"github.com/pkg/errors"
	"github.com/yozel/otrera/internal/log"
	"github.com/yozel/otrera/internal/objectstore"
	"github.com/yozel/otrera/internal/template"
)

func renderGoTemplate(tmpl string) (string, error) {
	logger := log.Log().With().Logger()
	var err error

	s, err := objectstore.NewObjectStore()
	if err != nil {
		panic(err)
	}

	template, err := template.New("hostTemplateString", string(tmpl), s)
	err = errors.Wrapf(err, "Can't parse hostTemplateString")
	if err != nil {
		return "nil", err
	}

	var b bytes.Buffer
	err = template.Execute(&b, map[string]interface{}{})
	if err != nil {
		logger.Fatal().Err(err).Msg("Can't execute template")
	}

	return b.String(), nil
}
