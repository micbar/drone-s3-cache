package plugin

import (
	"github.com/urfave/cli/v3"
)

// Flags returns all cli flags
func (p *Plugin) Flags() []cli.Flag {
	return []cli.Flag{
		// Cache information

		&cli.StringFlag{
			Name:        "mode",
			Usage:       "set plugin mode (rebuild,restore,flush)",
			Sources:     cli.EnvVars("PLUGIN_MODE"),
			Destination: &p.Settings.Mode,
			Value:       "rebuild",
		},
		&cli.StringFlag{
			Name:        "filename",
			Usage:       "filename for the cache archive",
			Sources:     cli.EnvVars("PLUGIN_FILENAME"),
			Destination: &p.Settings.Filename,
		},
		&cli.StringFlag{
			Name:        "root",
			Usage:       "storage root of cache files",
			Sources:     cli.EnvVars("PLUGIN_ROOT"),
			Destination: &p.Settings.Root,
		},
		&cli.StringFlag{
			Name:        "path",
			Usage:       "path to cache files relative to root",
			Sources:     cli.EnvVars("PLUGIN_PATH"),
			Destination: &p.Settings.Path,
		},
		&cli.StringFlag{
			Name:        "fallback-path",
			Usage:       "path to default cache files relative to the root",
			Sources:     cli.EnvVars("PLUGIN_FALLBACK_PATH"),
			Destination: &p.Settings.FallbackPath,
		},
		&cli.StringFlag{
			Name:        "flush-path",
			Usage:       "path to flushable cache files relative to the root",
			Sources:     cli.EnvVars("PLUGIN_FLUSH_PATH"),
			Destination: &p.Settings.FlushPath,
		},
		&cli.StringSliceFlag{
			Name:        "mount",
			Usage:       "directories to cache",
			Sources:     cli.EnvVars("PLUGIN_MOUNT"),
			Destination: &p.Settings.Mount,
		},
		&cli.IntFlag{
			Name:        "flush-age",
			Usage:       "flush cache files older than # days",
			Sources:     cli.EnvVars("PLUGIN_FLUSH_AGE"),
			Value:       30,
			Destination: &p.Settings.FlushAge,
		},

		// Cache information (deprecated)

		&cli.BoolFlag{
			Name:        "rebuild",
			Usage:       "rebuild the cache directories",
			Sources:     cli.EnvVars("PLUGIN_REBUILD"),
			Destination: &p.Settings.Rebuild,
		},
		&cli.BoolFlag{
			Name:        "restore",
			Usage:       "restore the cache directories",
			Sources:     cli.EnvVars("PLUGIN_RESTORE"),
			Destination: &p.Settings.Restore,
		},
		&cli.BoolFlag{
			Name:        "flush",
			Usage:       "flush the cache",
			Sources:     cli.EnvVars("PLUGIN_FLUSH"),
			Destination: &p.Settings.Flush,
		},

		// S3 information

		&cli.StringFlag{
			Name:        "endpoint",
			Usage:       "s3 endpoint",
			Sources:     cli.EnvVars("PLUGIN_SERVER", "PLUGIN_ENDPOINT", "CACHE_S3_ENDPOINT", "CACHE_S3_SERVER", "S3_ENDPOINT"),
			Destination: &p.Settings.S3Options.Endpoint,
		},
		&cli.StringFlag{
			Name:        "accelerated-endpoint",
			Usage:       "s3 accelerated endpoint",
			Sources:     cli.EnvVars("PLUGIN_ACCELERATED_ENDPOINT", "CACHE_S3_ACCELERATED_ENDPOINT"),
			Destination: &p.Settings.S3Options.AcceleratedEndpoint,
		},
		&cli.StringFlag{
			Name:        "access-key",
			Usage:       "s3 access key",
			Sources:     cli.EnvVars("PLUGIN_ACCESS_KEY", "CACHE_S3_ACCESS_KEY", "AWS_ACCESS_KEY_ID"),
			Destination: &p.Settings.S3Options.Access,
		},
		&cli.StringFlag{
			Name:        "secret-key",
			Usage:       "s3 secret key",
			Sources:     cli.EnvVars("PLUGIN_SECRET_KEY", "CACHE_S3_SECRET_KEY", "AWS_SECRET_ACCESS_KEY"),
			Destination: &p.Settings.S3Options.Secret,
		},
		&cli.StringFlag{
			Name:        "session-token",
			Usage:       "s3 session token",
			Sources:     cli.EnvVars("PLUGIN_SESSION_TOKEN", "CACHE_S3_SESSION_TOKEN", "AWS_SESSION_TOKEN"),
			Destination: &p.Settings.S3Options.Token,
		},
		&cli.StringFlag{
			Name:        "region",
			Usage:       "s3 region",
			Sources:     cli.EnvVars("PLUGIN_REGION", "CACHE_S3_REGION"),
			Destination: &p.Settings.S3Options.Region,
		},
		&cli.StringFlag{
			Name:        "file-credentials",
			Usage:       "path to s3 credentials file",
			Sources:     cli.EnvVars("PLUGIN_FILE_CREDENTIALS", "CACHE_FILE_CREDENTIALS", "AWS_SHARED_CREDENTIALS_FILE"),
			Destination: &p.Settings.S3Options.FileCredentials,
		},
		&cli.StringFlag{
			Name:        "profile",
			Usage:       "s3 profile name",
			Sources:     cli.EnvVars("PLUGIN_PROFILE", "CACHE_S3_PROFILE", "AWS_PROFILE"),
			Destination: &p.Settings.S3Options.Profile,
		},
	}
}
