// Package config handles CLI flags and configuration.
package config

import "flag"

// Version is the application version, can be overridden at build time.
var Version = "0.1.0"

// Config holds all configuration options for the application.
type Config struct {
	Version     *bool
	Help        *bool
	Output      *string
	Include     *string
	Exclude     *string
	Format      *string
	ShowTokens  *bool
	MaxFileSize *int
	MaxTokens   *int
}

// ParseFlags parses command-line flags and returns a Config.
func ParseFlags() *Config {
	cfg := &Config{
		Version:     flag.Bool("v", false, "print version"),
		Help:        flag.Bool("h", false, "show help"),
		Output:      flag.String("o", "", "output file (default stdout)"),
		Include:     flag.String("include", "", "comma-separated glob(s) to include (supports *, ?, [class])"),
		Exclude:     flag.String("exclude", "", "comma-separated glob(s) to exclude (supports *, ?, [class])"),
		Format:      flag.String("format", "markdown", "output format: markdown|json"),
		ShowTokens:  flag.Bool("tokens", false, "print estimated token count"),
		MaxFileSize: flag.Int("max-file-size", 16*1024, "per-file size limit in bytes before truncation"),
		MaxTokens:   flag.Int("max-tokens", 0, "stop when total estimated tokens reach this number (0 = no limit)"),
	}

	flag.Parse()
	return cfg
}

// Args returns the non-flag arguments.
func Args() []string {
	return flag.Args()
}
