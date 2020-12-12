package renderer

import (
	"bytes"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/yozel/otrera/gatherer/aws"
	"github.com/yozel/otrera/log"
	"github.com/yozel/otrera/objectstore"
	"github.com/yozel/otrera/template"
)

func renderGoTemplate(tmpl string) (string, error) {
	logger := log.Log().With().Logger()
	var err error

	configFile := "/Users/yasin.ozel/.aws/config"
	profiles, err := aws.ListProfiles(configFile)
	if err != nil {
		panic(err)
	}

	s, err := objectstore.NewObjectStore()
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(len(profiles))
	for _, profile := range profiles {
		go func(profile string) {
			defer wg.Done()
			logger.Info().Str("profile", profile).Msg("Processing profile")
			options := map[string]string{"profile": profile, "region": "eu-west-1"}
			labels := map[string]string{"profile": profile, "region": "eu-west-1"}
			err = s.Gather("AWS/EC2Instances", options, labels, 10*time.Minute)
			if err != nil {
				panic(err)
			}
			err = s.Gather("AWS/EC2Images", options, labels, 10*time.Minute)
			if err != nil {
				panic(err)
			}
			logger.Info().Str("profile", profile).Msg("Done processing profile")
		}(profile)
	}
	wg.Wait()

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
