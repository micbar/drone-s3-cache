// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package plugin

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/url"
	"os"
	pathutil "path"
	"strings"
	"time"

	"github.com/drone/drone-cache-lib/archive/util"
	"github.com/drone/drone-cache-lib/cache"
	"github.com/drone/drone-cache-lib/storage"
	"github.com/micbar/drone-s3-cache/storage/s3"
)

// Validate handles the settings validation of the plugin.
func (p *Plugin) Validate() error {
	if err := p.validateMode(); err != nil {
		return err
	}
	return p.validateS3()
}

func (p *Plugin) validateMode() error {
	// Validate the mode
	mode := p.Settings.Mode
	hasMode := p.Settings.Rebuild || p.Settings.Restore || p.Settings.Flush
	if mode == "" {
		logrus.WithFields(logrus.Fields{
			"rebuild": p.Settings.Rebuild,
			"restore": p.Settings.Restore,
			"flush":   p.Settings.Flush,
		}).Info("mode specified using boolean config")

		if !hasMode {
			return fmt.Errorf("no mode specified")
		}
		if multipleModesSpecified(p.Settings.Rebuild, p.Settings.Restore, p.Settings.Flush) {
			return fmt.Errorf("multiple modes specified")
		}

		if p.Settings.Rebuild {
			mode = rebuildMode
		} else if p.Settings.Restore {
			mode = restoreMode
		} else {
			mode = flushMode
		}
	} else {
		if hasMode {
			return fmt.Errorf("mode specified multiple ways")
		}

		if mode != rebuildMode && mode != restoreMode && mode != flushMode {
			return fmt.Errorf("invalid mode %s specified", mode)
		}
	}

	logrus.WithField("mode", mode).Info("using mode")
	p.Settings.Mode = mode

	if p.Settings.Filename == "" {
		logrus.Debug("using default filename")
		p.Settings.Filename = "archive.tar"
	}
	logrus.WithField("filename", p.Settings.Filename).Debug("using filename")

	// Validate mode settings
	if mode != flushMode {
		if p.Settings.Path == "" {
			logrus.WithFields(logrus.Fields{
				"repo.owner":  p.Metadata.Repository.Owner,
				"repo.name":   p.Metadata.Repository.Name,
				"repo.branch": p.Metadata.Repository.DefaultBranch,
			}).Debug("creating default path")
			p.Settings.Path = fmt.Sprintf(
				"%s/%s/%s",
				p.Metadata.Repository.Owner,
				p.Metadata.Repository.Name,
				p.Metadata.Repository.DefaultBranch,
			)
		}
		logrus.WithField("path", p.Settings.Path).Debug("using path")

		if mode == rebuildMode {
			mount := p.Settings.Mount
			if len(mount) == 0 {
				return fmt.Errorf("cache not specified")
			}
			p.Settings.mount = mount
		} else {
			if p.Settings.FallbackPath == "" {
				logrus.WithFields(logrus.Fields{
					"repo.owner":  p.Metadata.Repository.Owner,
					"repo.name":   p.Metadata.Repository.Name,
					"repo.branch": p.Metadata.Repository.DefaultBranch,
				}).Debug("creating default fallback path")
				p.Settings.FallbackPath = fmt.Sprintf(
					"%s/%s/%s",
					p.Metadata.Repository.Owner,
					p.Metadata.Repository.Name,
					p.Metadata.Repository.DefaultBranch,
				)
			}
			logrus.WithField("path", p.Settings.FallbackPath).Debug("using path as fallback")
		}
	} else {
		if p.Settings.FlushPath == "" {
			logrus.WithFields(logrus.Fields{
				"repo.owner": p.Metadata.Repository.Owner,
				"repo.name":  p.Metadata.Repository.Name,
			}).Debug("creating default flush path")
			p.Settings.FlushPath = fmt.Sprintf(
				"%s/%s",
				p.Metadata.Repository.Owner,
				p.Metadata.Repository.Name,
			)
		}
		logrus.WithField("path", p.Settings.FlushPath).Debug("using path when flushing")
	}

	return nil
}

func (p *Plugin) validateS3() error {
	// Validate the endpoint
	endpoint := p.Settings.S3Options.Endpoint
	isAWS := false
	bucket := ""
	region := ""

	if endpoint == "" {
		endpoint = awsEndpoint
	}

	s3url, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("could not parse endpoint %s", endpoint)
	}

	// Check if additional information is encoded in the endpoint
	if h := s3url.Hostname(); strings.HasSuffix(h, awsDomain) {
		isAWS = true

		// Sub-domains can specify region and bucket
		d := strings.Split(h, ".")
		s3sub := 0

		switch len(d) {
		case 5:
			// Virtual hosted style access
			// https://bucket-name.s3.Region.amazonaws.com/
			logrus.WithField("host", h).Debug("using virtual host style access")
			bucket = d[0]
			s3sub = 1
			region = d[2]
		case 4:
			// Path-style access
			// https://s3.Region.amazonaws.com/bucket-name
			logrus.WithField("host", h).Debug("using path style access")
			bucket = s3url.Path
			s3sub = 0
			region = d[1]
		case 3:
			// Just default url https://s3.Region.amazonaws.com
		default:
			return fmt.Errorf("unknown aws domain for url %s", endpoint)
		}

		if d[s3sub] != "s3" {
			return fmt.Errorf("unknown aws domain for url %s", endpoint)
		}

		// Normalize endpoint
		endpoint = awsEndpoint
		s3url, _ = url.Parse(endpoint)
	}

	// Check for s3 scheme
	if s3url.Scheme == "s3" {
		logrus.WithField("endpoint", endpoint).Debug("using s3 url")
		bucket = s3url.Hostname()

		// Normalize endpoint
		endpoint = awsEndpoint
		s3url, _ = url.Parse(endpoint)
	}

	var useSSL bool
	switch s3url.Scheme {
	case "https":
		endpoint = endpoint[8:]
		useSSL = true
	case "http":
		endpoint = endpoint[7:]
		useSSL = false
	default:
		return fmt.Errorf("unknown scheme for endpoint %s", endpoint)
	}

	if bucket != "" {
		logrus.WithField("bucket", bucket).Info("bucket found in S3 endpoint")
		if p.Settings.Root != "" {
			return fmt.Errorf("bucket %s already specified in endpoint remove from root", bucket)
		}
		p.Settings.Root = bucket
	}

	if region != "" {
		logrus.WithField("region", region).Info("region found in S3 endpoint")
		if p.Settings.S3Options.Region != "" {
			return fmt.Errorf("region %s already specified in endpoint remove from config", region)
		}
		p.Settings.S3Options.Region = region
	}
	s3Opts := p.Settings.S3Options

	if (s3Opts.Access != "" || s3Opts.Secret != "") && s3Opts.FileCredentials != "" {
		return fmt.Errorf("only one credentials method should be used. Use either access-key and secret-key OR the credentials file")
	}

	if s3Opts.FileCredentials != "" {
		if _, err := os.Stat(s3Opts.FileCredentials); os.IsNotExist(err) {
			return fmt.Errorf("file %s does not exist", s3Opts.FileCredentials)
		}
	}

	logrus.WithFields(logrus.Fields{
		"endpoint": endpoint,
		"use-ssl":  useSSL,
	}).Info("using S3 endpoint")
	p.Settings.S3Options.Endpoint = endpoint
	p.Settings.S3Options.UseSSL = useSSL

	if isAWS && p.Settings.Root == "" {
		return fmt.Errorf("no aws bucket specified in root or endpoint")
	}

	return nil
}

// Execute provides the implementation of the plugin.
func (p *Plugin) Execute(ctx context.Context) error {
	err := p.Validate()
	if err != nil {
		return err
	}
	at, err := util.FromFilename(p.Settings.Filename)
	if err != nil {
		return err
	}

	st, err := s3.New(&p.Settings.S3Options)
	if err != nil {
		return err
	}

	c := cache.New(st, at)

	if p.Settings.Mode == rebuildMode {
		path := cleanPath(p.Settings.Root, p.Settings.Path, p.Settings.Filename)
		logrus.WithFields(logrus.Fields{
			"path": path,
		}).Info("rebuilding cache")
		err = c.Rebuild(p.Settings.mount, path)

		if err == nil {
			logrus.Infof("cache rebuilt")
		}
	} else if p.Settings.Mode == restoreMode {
		path := cleanPath(p.Settings.Root, p.Settings.Path, p.Settings.Filename)
		fallbackPath := cleanPath(p.Settings.Root, p.Settings.FallbackPath, p.Settings.Filename)

		logrus.WithFields(logrus.Fields{
			"path":     path,
			"fallback": fallbackPath,
		}).Info("restoring cache")
		err = c.Restore(path, fallbackPath)

		if err == nil {
			logrus.Info("cache restored")
		}
	} else /* p.Settings.Mode == flushMode */ {
		flushPath := cleanPath(p.Settings.Root, p.Settings.FlushPath)

		logrus.WithFields(logrus.Fields{
			"path":    flushPath,
			"max-age": p.Settings.FlushAge,
		}).Info("flushing cache")
		f := cache.NewFlusher(st, genIsExpired(int(p.Settings.FlushAge)))
		err = f.Flush(flushPath)

		if err == nil {
			logrus.Info("Cache flushed")
		}
	}

	return err
}

func genIsExpired(age int) cache.DirtyFunc {
	return func(file storage.FileEntry) bool {
		// Check if older than "age" days
		return file.LastModified.Before(time.Now().AddDate(0, 0, age*-1))
	}
}

func cleanPath(paths ...string) string {
	return pathutil.Clean(pathutil.Join(paths...))
}

func multipleModesSpecified(bools ...bool) bool {
	var b bool
	for _, v := range bools {
		if b && b == v {
			return true
		}

		if v {
			b = true
		}
	}

	return false
}
