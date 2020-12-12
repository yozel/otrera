package main

import (
	"github.com/rs/zerolog"
	"github.com/yozel/otrera/cmd"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	cmd.Execute()
}
