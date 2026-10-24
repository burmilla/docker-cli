package system

import (
	"fmt"
	"strings"

	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/command/inspect"
	"github.com/docker/docker/api/types"
	apiclient "github.com/docker/docker/client"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

type inspectOptions struct {
	format      string
	inspectType string
	size        bool
	ids         []string
}

// NewInspectCommand creates a new cobra.Command for `docker inspect`
func NewInspectCommand(dockerCli *command.DockerCli) *cobra.Command {
	var opts inspectOptions

	cmd := &cobra.Command{
		Use:   "inspect [OPTIONS] NAME|ID [NAME|ID...]",
		Short: "Return low-level information on Docker objects",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ids = args
			return runInspect(dockerCli, opts)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.format, "format", "f", "", "Format the output using the given Go template")
	flags.StringVar(&opts.inspectType, "type", "", "Return JSON for specified type")
	flags.BoolVarP(&opts.size, "size", "s", false, "Display total file sizes if the type is container")

	return cmd
}

func runInspect(dockerCli *command.DockerCli, opts inspectOptions) error {
	var elementSearcher inspect.GetRefFunc
	switch opts.inspectType {
	case "", "container", "image", "network", "volume":
		elementSearcher = inspectAll(context.Background(), dockerCli, opts.size, opts.inspectType)
	default:
		return errors.Errorf("%q is not a valid value for --type", opts.inspectType)
	}
	return inspect.Inspect(dockerCli.Out(), opts.ids, opts.format, elementSearcher)
}

func inspectContainers(ctx context.Context, dockerCli *command.DockerCli, getSize bool) inspect.GetRefFunc {
	return func(ref string) (interface{}, []byte, error) {
		return dockerCli.Client().ContainerInspectWithRaw(ctx, ref, getSize)
	}
}

func inspectImages(ctx context.Context, dockerCli *command.DockerCli) inspect.GetRefFunc {
	return func(ref string) (interface{}, []byte, error) {
		return dockerCli.Client().ImageInspectWithRaw(ctx, ref)
	}
}

func inspectNetwork(ctx context.Context, dockerCli *command.DockerCli) inspect.GetRefFunc {
	return func(ref string) (interface{}, []byte, error) {
		return dockerCli.Client().NetworkInspectWithRaw(ctx, ref, types.NetworkInspectOptions{})
	}
}

func inspectVolume(ctx context.Context, dockerCli *command.DockerCli) inspect.GetRefFunc {
	return func(ref string) (interface{}, []byte, error) {
		return dockerCli.Client().VolumeInspectWithRaw(ctx, ref)
	}
}

func inspectAll(ctx context.Context, dockerCli *command.DockerCli, getSize bool, typeConstraint string) inspect.GetRefFunc {
	var inspectAutodetect = []struct {
		objectType      string
		isSizeSupported bool
		objectInspector func(string) (interface{}, []byte, error)
	}{
		{
			objectType:      "container",
			isSizeSupported: true,
			objectInspector: inspectContainers(ctx, dockerCli, getSize),
		},
		{
			objectType:      "image",
			objectInspector: inspectImages(ctx, dockerCli),
		},
		{
			objectType:      "network",
			objectInspector: inspectNetwork(ctx, dockerCli),
		},
		{
			objectType:      "volume",
			objectInspector: inspectVolume(ctx, dockerCli),
		},
	}

	isErrNotSupported := func(err error) bool {
		return strings.Contains(err.Error(), "not supported")
	}

	return func(ref string) (interface{}, []byte, error) {
		for _, inspectData := range inspectAutodetect {
			if typeConstraint != "" && inspectData.objectType != typeConstraint {
				continue
			}
			v, raw, err := inspectData.objectInspector(ref)
			if err != nil {
				if typeConstraint == "" && (apiclient.IsErrNotFound(err) || isErrNotSupported(err)) {
					continue
				}
				return v, raw, err
			}
			if getSize && !inspectData.isSizeSupported {
				fmt.Fprintf(dockerCli.Err(), "WARNING: --size ignored for %s\n", inspectData.objectType)
			}
			return v, raw, err
		}
		return nil, nil, errors.Errorf("Error: No such object: %s", ref)
	}
}
