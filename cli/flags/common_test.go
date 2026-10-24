package flags

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestCommonOptionsInstallFlags(t *testing.T) {
	flags := pflag.NewFlagSet("testing", pflag.ContinueOnError)
	opts := NewCommonOptions()
	opts.InstallFlags(flags)

	err := flags.Parse([]string{"-H", "unix:///run/system-docker.sock"})
	assert.NoError(t, err)
	assert.Equal(t, []string{"unix:///run/system-docker.sock"}, opts.Hosts)
}
