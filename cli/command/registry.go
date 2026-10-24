package command

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/types"
	registrytypes "github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/registry"
)

// ElectAuthServer returns the default registry to use (by asking the daemon)
func ElectAuthServer(ctx context.Context, cli Cli) string {
	// The daemon `/info` endpoint informs us of the default registry being
	// used. This is essential in cross-platforms environment, where for
	// example a Linux client might be interacting with a Windows daemon, hence
	// the default registry URL might be Windows specific.
	serverAddress := registry.IndexServer
	if info, err := cli.Client().Info(ctx); err != nil {
		fmt.Fprintf(cli.Err(), "Warning: failed to get default registry endpoint from daemon (%v). Using system default: %s\n", err, serverAddress)
	} else if info.IndexServerAddress == "" {
		fmt.Fprintf(cli.Err(), "Warning: Empty registry endpoint from daemon. Using system default: %s\n", serverAddress)
	} else {
		serverAddress = info.IndexServerAddress
	}
	return serverAddress
}

// EncodeAuthToBase64 serializes the auth configuration as JSON base64 payload
func EncodeAuthToBase64(authConfig types.AuthConfig) (string, error) {
	buf, err := json.Marshal(authConfig)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(buf), nil
}

// ResolveAuthConfig is like registry.ResolveAuthConfig, but if using the
// default index, it uses the default index name for the daemon's platform,
// not the client's platform.
func ResolveAuthConfig(ctx context.Context, cli Cli, index *registrytypes.IndexInfo) types.AuthConfig {
	configKey := index.Name
	if index.Official {
		configKey = ElectAuthServer(ctx, cli)
	}

	a, _ := cli.CredentialsStore(configKey).Get(configKey)
	return a
}
