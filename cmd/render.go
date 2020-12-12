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
	"github.com/yozel/otrera/internal/renderer"

	"github.com/spf13/cobra"
)

var (
	flagTemplatePath    string
	flagPrependFilePath string
	flagAppendFilePath  string
)

// renderCmd represents the render command
var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Render a template with AWS data",
	Run: handleErrors(func(cmd *cobra.Command, args []string) error {
		// logger := log.Log().With().Logger()
		// loggerDebug := log.Log().With().
		// 	Str("cobra cmd", "render").
		// 	Str("flagPrependFilePath", flagPrependFilePath).
		// 	Str("flagAppendFilePath", flagAppendFilePath).Logger()

		if flagPrependFilePath != "" {
			if err := renderer.NewRenderableWithPath(renderer.Text, flagPrependFilePath).Render(); err != nil {
				return err
			}
		}

		if err := renderer.NewRenderableWithPath(renderer.GoTemplate, flagTemplatePath).Render(); err != nil {
			return err
		}

		if flagAppendFilePath != "" {
			if err := renderer.NewRenderableWithPath(renderer.Text, flagAppendFilePath).Render(); err != nil {
				return err
			}
		}

		return nil
	}),
}

func init() {
	rootCmd.AddCommand(renderCmd)
	renderCmd.PersistentFlags().StringVarP(&flagTemplatePath, "template", "t", "", "Template to render")
	renderCmd.PersistentFlags().StringVar(&flagPrependFilePath, "prepend-file", "", "File to prepend before rendered file")
	renderCmd.PersistentFlags().StringVar(&flagAppendFilePath, "append-file", "", "File to append after renderes file")
	renderCmd.MarkPersistentFlagRequired("template")
	renderCmd.MarkPersistentFlagFilename("template")
	renderCmd.MarkPersistentFlagFilename("prepend-file")
	renderCmd.MarkPersistentFlagFilename("append-file")
}
