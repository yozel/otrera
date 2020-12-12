/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"github.com/yozel/otrera/template"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/yozel/otrera/gatherer/aws"
	"github.com/yozel/otrera/objectstore"
)

var (
	flagTemplatePath string
)

// renderCmd represents the render command
var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Render a template with AWS data",
	Run: func(cmd *cobra.Command, args []string) {
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

		hostTemplateString, err := ioutil.ReadFile(flagTemplatePath)

		template, err := template.New("hostTemplateString", string(hostTemplateString), s)
		err = errors.Wrapf(err, "Can't parse hostTemplateString")
		if err != nil {
			panic(err)
		}

		var b bytes.Buffer
		err = template.Execute(&b, map[string]interface{}{})
		err = errors.Wrapf(err, "Can't execute template")
		if err != nil {
			log.Println(err)
		}

		fmt.Printf("%s", b.String())
	},
}

func init() {
	rootCmd.AddCommand(renderCmd)
	renderCmd.PersistentFlags().StringVarP(&flagTemplatePath, "template", "t", "", "Template to render")
	renderCmd.MarkPersistentFlagRequired("template")
	renderCmd.MarkPersistentFlagFilename("template")
}
