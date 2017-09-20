package restapi

import (
	"fmt"

	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	defaultAddress = ":8080"
)

// Config defines the config options for the API server.
type Config struct {
	Address           string
	Debug             bool
	InsecureHTTP      bool
	AuthDisabled      bool
	TLSCertFile       string
	TLSKeyFile        string
	WellKnownDisabled bool
	TokenURL          string
}

// Parse parses configuration from commandline flags.
func (c *Config) Parse() error {
	kingpin.Flag("debug", "Enable debug logging and pprof metrics.").BoolVar(&c.Debug)
	kingpin.Flag("address", "Address to listen on, e.g. :8080 or 0.0.0.0:8080.").
		Default(defaultAddress).StringVar(&c.Address)
	kingpin.Flag("insecure-http", "Service only HTTP.").BoolVar(&c.InsecureHTTP)
	kingpin.Flag("tls-cert-file", "Path to TLS Cert file used when serving HTTPS.").
		StringVar(&c.TLSCertFile)
	kingpin.Flag("tls-key-file", "Path to TLS Key file used when serving HTTPS.").
		StringVar(&c.TLSKeyFile)
	kingpin.Flag("disable-well-known", "Disable automatic /.well-known resources.").
		BoolVar(&c.WellKnownDisabled)
	kingpin.Flag("disable-auth", "Disable auth for all resources.").
		BoolVar(&c.AuthDisabled)
	kingpin.Flag("token-url", "Set TokenURL used to validate oauth2 tokens.").
		StringVar(&c.TokenURL)
	kingpin.Parse()

	if !c.InsecureHTTP && (c.TLSCertFile == "" || c.TLSKeyFile == "") {
		return fmt.Errorf("'--tls-cert-file' and '--tls-key-file' must be specified when '--insecure-http=false'")
	}

	return nil
}

// vim: ft=go
