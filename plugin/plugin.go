// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package plugin

import (
	"codeberg.org/woodpecker-plugins/go-plugin"
	"github.com/micbar/drone-s3-cache/storage/s3"
)

// Plugin implements drone.Plugin to provide the plugin implementation.
type Plugin struct {
	*plugin.Plugin
	Settings Settings
}

// Settings for the plugin.
type Settings struct {
	Mode         string
	Root         string
	Filename     string
	Path         string
	FallbackPath string
	FlushPath    string
	FlushAge     int64
	Mount        []string
	Restore      bool // DEPRECATED
	Rebuild      bool // DEPRECATED
	Flush        bool // DEPRECATED

	S3Options s3.Options
	mount     []string
}

const (
	restoreMode = "restore"
	rebuildMode = "rebuild"
	flushMode   = "flush"

	awsDomain   = "amazonaws.com"
	awsEndpoint = "https://s3." + awsDomain
)
