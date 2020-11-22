package main

import (
	"bytes"
	"fmt"
	"log"
	"sync"

	"time"

	"github.com/pkg/errors"
	"github.com/yozel/otrera/gatherer/aws"
	"github.com/yozel/otrera/objectstore"
	"github.com/yozel/otrera/template"
)

var hostTemplateString string = `#### START OF THE TEMPLATE ####
{{ range $object := getall "AWS/EC2Instances/*" }}
{{- with $object.Data }}
{{- if .PrivateIpAddress }}
Host {{ $object.Labels.profile }}_{{ .Tags | gjson "#(Key==\"Name\").Value" }}_{{ .PrivateIpAddress | replace "." "-" }}
	Hostname {{ .PrivateIpAddress }}
	{{ if .KeyName }}IdentityFile ~/.ssh/subdir/{{ .KeyName }}.pem{{ end }}
	# ImageId: {{ .ImageId }}
{{ end -}}
{{ end -}}
{{ end }}
#### END OF THE TEMPLATE ####
`

func main() {
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
			log.Printf("Processing profile %s\n", profile)
			options := map[string]string{"profile": profile, "region": "eu-west-1"}
			labels := map[string]string{"profile": profile, "region": "eu-west-1"}
			err = s.Gather("AWS/EC2Instances", options, labels, 10*time.Minute)
			if err != nil {
				panic(err)
			}
			log.Printf("Done processing profile %s\n", profile)
		}(profile)
	}
	wg.Wait()

	template, err := template.New("hostTemplateString", hostTemplateString, s)
	err = errors.Wrapf(err, "Can't parse hostTemplateString")
	if err != nil {
		panic(err)
	}

	var b bytes.Buffer
	err = template.Execute(&b, map[string]interface{}{
		// "ObjectStore": s,
	})
	err = errors.Wrapf(err, "Can't execute template")
	if err != nil {
		log.Println(err)
	}

	fmt.Printf("%s", b.String())

	// s.GetAll("AWS/EC2Instances/i-09cb7d4b65abf60d2", nil)

}
