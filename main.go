package main

import (
	"bytes"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/yozel/otrera/gatherer"
	"github.com/yozel/otrera/gatherer/aws"
	"github.com/yozel/otrera/template"

	"github.com/pkg/errors"
)

var hostTemplateString string = `#### START OF THE TEMPLATE ####

{{ range $_, $instance := $.Instances }}
{{- if $instance.PrivateIpAddress }}
Host {{ $.Profile }}_{{ $instance.Tags | gjson "#(Key==\"Name\").Value" }}_{{ $instance.PrivateIpAddress | replace "." "-" }}
	Hostname {{ $instance.PrivateIpAddress }}
	{{ if $instance.KeyName }}IdentityFile ~/.ssh/subdir/{{ $instance.KeyName }}.pem{{ end }}
	# ImageId: {{ $instance.ImageId }}
{{- end }}
{{ end }}
#### END OF THE TEMPLATE ####
`

func main() {
	template, err := template.New(hostTemplateString)
	err = errors.Wrapf(err, "Can't parse hostTemplateString")
	if err != nil {
		panic(err)
	}

	configFile := "/Users/yasin.ozel/.aws/config"
	profiles, err := aws.ListProfiles(configFile)
	if err != nil {
		panic(err)
	}
	g := gatherer.New("/tmp/.otrera.cache")

	var wg sync.WaitGroup
	wg.Add(len(profiles))

	results := make(chan string)
	for _, profile := range profiles {
		go func(profile string) {
			defer wg.Done()
			log.Printf("Processing profile %s\n", profile)

			c, err := g.Gather("EC2Instances", map[string]string{"profile": profile, "region": "eu-west-1"}, 10*time.Minute)
			if err != nil {
				panic(err)
			}
			var b bytes.Buffer
			err = template.Execute(&b, map[string]interface{}{
				"Profile":   profile,
				"Instances": c.Data,
			})
			err = errors.Wrapf(err, "Can't execute template")
			if err != nil {
				log.Println(err)
			}

			log.Printf("Done processing profile %s\n", profile)
			results <- b.String()
		}(profile)
	}

	for i := 0; i < len(profiles); i++ {
		fmt.Printf("%s", <-results)
	}
	close(results)
	wg.Wait()
}
