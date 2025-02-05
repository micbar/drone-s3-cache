// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

// DO NOT MODIFY THIS FILE DIRECTLY

package main

import (
	gplugin "codeberg.org/woodpecker-plugins/go-plugin"
	"github.com/joho/godotenv"
	"github.com/micbar/drone-s3-cache/plugin"
	"github.com/rs/zerolog/log"
	"os"
)

var version = "unknown"

func main() {
	if envFile, set := os.LookupEnv("PLUGIN_ENV_FILE"); set {
		if err := godotenv.Overload(envFile); err != nil {
			log.Error().Err(err).Msg("couldn't load env file")
		}
	}
	p := &plugin.Plugin{
		Settings: plugin.Settings{},
	}

	p.Plugin = gplugin.New(gplugin.Options{
		Name:        "s3-cache",
		Description: "Plugin to cache build artifacts in S3",
		Flags:       p.Flags(),
		Execute:     p.Execute,
	})

	p.Run()
}
