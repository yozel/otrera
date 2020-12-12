package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yozel/otrera/internal/renderer"
)

var (
	flagRenderers []string
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

		for _, r := range flagRenderers {
			parts := strings.SplitN(r, ":", 3)

			rtype := parts[0]
			source := parts[1]

			var rendererType renderer.RenderableType

			switch rtype {
			case "static":
				rendererType = renderer.Static
			case "gotmpl":
				rendererType = renderer.GoTemplate
			default:
				return fmt.Errorf("Unknown Renderable Type: %s", rtype)
			}

			var err error

			switch source {
			case "file":
				err = renderer.NewRenderableWithPath(rendererType, parts[2]).Render()
			case "literal":
				err = renderer.NewRenderableWithContent(rendererType, parts[2]).Render()
			}
			if err != nil {
				return err
			}
		}

		return nil
	}),
}

func init() {
	rootCmd.AddCommand(renderCmd)
	renderCmd.PersistentFlags().StringArrayVarP(&flagRenderers, "renderer", "r", []string{}, "")
	renderCmd.MarkPersistentFlagRequired("renderer")
}
