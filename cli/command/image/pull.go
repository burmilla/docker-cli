package image

import (
	"fmt"

	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli/command"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/registry"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

type pullOptions struct {
	remote string
	all    bool
}

// NewPullCommand creates a new `docker pull` command
func NewPullCommand(dockerCli command.Cli) *cobra.Command {
	var opts pullOptions

	cmd := &cobra.Command{
		Use:   "pull [OPTIONS] NAME[:TAG|@DIGEST]",
		Short: "Pull an image or a repository from a registry",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.remote = args[0]
			return runPull(dockerCli, opts)
		},
	}

	flags := cmd.Flags()

	flags.BoolVarP(&opts.all, "all-tags", "a", false, "Download all tagged images in the repository")

	return cmd
}

func runPull(dockerCli command.Cli, opts pullOptions) error {
	distributionRef, err := reference.ParseNormalizedNamed(opts.remote)
	if err != nil {
		return err
	}
	if opts.all && !reference.IsNameOnly(distributionRef) {
		return errors.New("tag can't be used with --all-tags/-a")
	}

	if !opts.all && reference.IsNameOnly(distributionRef) {
		distributionRef = reference.TagNameOnly(distributionRef)
		if tagged, ok := distributionRef.(reference.Tagged); ok {
			fmt.Fprintf(dockerCli.Out(), "Using default tag: %s\n", tagged.Tag())
		}
	}

	// Resolve the Repository name from fqn to RepositoryInfo
	repoInfo, err := registry.ParseRepositoryInfo(distributionRef)
	if err != nil {
		return err
	}

	ctx := context.Background()

	encodedAuth, err := command.EncodeAuthToBase64(command.ResolveAuthConfig(ctx, dockerCli, repoInfo.Index))
	if err != nil {
		return err
	}
	options := types.ImagePullOptions{
		RegistryAuth: encodedAuth,
		All:          opts.all,
	}

	responseBody, err := dockerCli.Client().ImagePull(ctx, reference.FamiliarString(distributionRef), options)
	if err != nil {
		return err
	}
	defer responseBody.Close()

	return jsonmessage.DisplayJSONMessagesToStream(responseBody, dockerCli.Out(), nil)
}
