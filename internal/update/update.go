package update

import (
	"fmt"
	"sync"
	"time"

	"github.com/yozel/otrera/internal/gatherer"
	"github.com/yozel/otrera/internal/gatherer/aws"
	"github.com/yozel/otrera/internal/log"
	"github.com/yozel/otrera/internal/objectstore"
)

func Update() error {
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

			g := gatherer.New(map[string]func(options map[string]string) ([]gatherer.RawObjectInterface, error){
				"AWS/EC2Instances": aws.DescribeEC2Instances,
				"AWS/EC2Images":    aws.DescribeEC2Images,
			})

			err = s.Clear()
			if err != nil {
				panic(err)
			}

			populateObjectStore := func(key string, o map[string]string, l map[string]string) error {
				c, err := g.Gather(key, o)
				if err != nil {
					return err // TODO: wrap error
				}
				for _, object := range c {
					s.Set(
						fmt.Sprintf("%s/%s", key, object.Name()),
						l,
						time.Now().UTC(),
						object.Content(),
					)
				}
				return nil
			}

			if err = populateObjectStore("AWS/EC2Instances", options, labels); err != nil {
				panic(err)
			}
			if err = populateObjectStore("AWS/EC2Images", options, labels); err != nil {
				panic(err)
			}

			logger.Info().Str("profile", profile).Msg("Done processing profile")
		}(profile)
	}
	wg.Wait()
	return nil
}
