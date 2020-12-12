package update

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/yozel/otrera/internal/gatherer/aws"
	"github.com/yozel/otrera/internal/log"
	"github.com/yozel/otrera/internal/objectstore"
	"github.com/yozel/otrera/internal/types"
)

var logger zerolog.Logger

func init() {
	logger = log.Log().With().Logger()
}

func Update() error {
	h, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configFile := h + "/.aws/config"
	profiles, err := aws.ListProfiles(configFile)
	if err != nil {
		return err
	}

	s, err := objectstore.NewObjectStore()
	if err != nil {
		return err
	}
	err = s.Clear()
	if err != nil {
		return err
	}

	descriptors := map[string]func(options map[string]string) ([]types.RawObjectInterface, error){
		"AWS/EC2Instances": aws.DescribeEC2Instances,
		"AWS/EC2Images":    aws.DescribeEC2Images,
	}

	var wg sync.WaitGroup
	for _, profile := range profiles {
		options := map[string]string{"profile": profile, "region": "eu-west-1"}
		labels := map[string]string{"profile": profile, "region": "eu-west-1"}
		for dn, d := range descriptors {
			wg.Add(1)
			go func(
				profile string,
				descritorName string,
				descriptor func(options map[string]string) ([]types.RawObjectInterface, error),
				options map[string]string,
				labels map[string]string,
			) {
				sublogger := logger.With().Str("profile", profile).Str("descriptor", descritorName).Logger()
				sublogger.Info().Msg("Processing started")
				defer wg.Done()
				c, err := descriptor(options)
				if err != nil {
					sublogger.Error().Err(err).Str("options", fmt.Sprintf("%+v", options)).Msg("Failed to add to objectstore")
					return
				}
				for _, object := range c {
					objName := fmt.Sprintf("%s/%s", descritorName, object.Name())
					err := s.Set(objName, labels, time.Now().UTC(), object.Content())
					if err != nil {
						sublogger.Error().Err(err).Str("objName", objName).Msg("Failed to add to objectstore")
						continue
					}
				}
				sublogger.Info().Msg("Processing done")
			}(profile, dn, d, options, labels)
		}
	}
	wg.Wait()
	return nil
}
