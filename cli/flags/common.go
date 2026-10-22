package flags

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/docker/cli/opts"
	"github.com/spf13/pflag"
)

// CommonOptions are options common to both the client and the daemon.
type CommonOptions struct {
	Debug    bool
	Hosts    []string
	LogLevel string
}

// NewCommonOptions returns a new CommonOptions
func NewCommonOptions() *CommonOptions {
	return &CommonOptions{}
}

// InstallFlags adds flags for the common options on the FlagSet
func (commonOpts *CommonOptions) InstallFlags(flags *pflag.FlagSet) {
	flags.BoolVarP(&commonOpts.Debug, "debug", "D", false, "Enable debug mode")
	flags.StringVarP(&commonOpts.LogLevel, "log-level", "l", "info", `Set the logging level ("debug"|"info"|"warn"|"error"|"fatal")`)
	hostOpt := opts.NewNamedListOptsRef("hosts", &commonOpts.Hosts, opts.ValidateHost)
	flags.VarP(hostOpt, "host", "H", "Daemon socket(s) to connect to")
}

// SetLogLevel sets the logrus logging level
func SetLogLevel(logLevel string) {
	if logLevel != "" {
		lvl, err := logrus.ParseLevel(logLevel)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to parse logging level: %s\n", logLevel)
			os.Exit(1)
		}
		logrus.SetLevel(lvl)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}
