package system

import (
	"fmt"
	"io"
	"strings"

	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/debug"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/templates"
	"github.com/docker/go-units"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

type infoOptions struct {
	format string
}

// NewInfoCommand creates a new cobra.Command for `docker info`
func NewInfoCommand(dockerCli *command.DockerCli) *cobra.Command {
	var opts infoOptions

	cmd := &cobra.Command{
		Use:   "info [OPTIONS]",
		Short: "Display system-wide information",
		Args:  cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInfo(dockerCli, &opts)
		},
	}

	flags := cmd.Flags()

	flags.StringVarP(&opts.format, "format", "f", "", "Format the output using the given Go template")

	return cmd
}

func runInfo(dockerCli *command.DockerCli, opts *infoOptions) error {
	ctx := context.Background()
	info, err := dockerCli.Client().Info(ctx)
	if err != nil {
		return err
	}
	if opts.format == "" {
		return prettyPrintInfo(dockerCli, info)
	}
	return formatInfo(dockerCli, info, opts.format)
}

// nolint: gocyclo
func prettyPrintInfo(dockerCli *command.DockerCli, info types.Info) error {
	fmt.Fprintf(dockerCli.Out(), "Containers: %d\n", info.Containers)
	fmt.Fprintf(dockerCli.Out(), " Running: %d\n", info.ContainersRunning)
	fmt.Fprintf(dockerCli.Out(), " Paused: %d\n", info.ContainersPaused)
	fmt.Fprintf(dockerCli.Out(), " Stopped: %d\n", info.ContainersStopped)
	fmt.Fprintf(dockerCli.Out(), "Images: %d\n", info.Images)
	fprintfIfNotEmpty(dockerCli.Out(), "Server Version: %s\n", info.ServerVersion)
	fprintfIfNotEmpty(dockerCli.Out(), "Storage Driver: %s\n", info.Driver)
	if info.DriverStatus != nil {
		for _, pair := range info.DriverStatus {
			fmt.Fprintf(dockerCli.Out(), " %s: %s\n", pair[0], pair[1])
		}

	}
	if info.SystemStatus != nil {
		for _, pair := range info.SystemStatus {
			fmt.Fprintf(dockerCli.Out(), "%s: %s\n", pair[0], pair[1])
		}
	}
	fprintfIfNotEmpty(dockerCli.Out(), "Logging Driver: %s\n", info.LoggingDriver)

	fmt.Fprintf(dockerCli.Out(), " Log:")
	fmt.Fprintf(dockerCli.Out(), " %s", strings.Join(info.Plugins.Log, " "))
	fmt.Fprintf(dockerCli.Out(), "\n")

	if info.OSType == "linux" {
		for _, ci := range []struct {
			Name   string
			Commit types.Commit
		}{
			{"containerd", info.ContainerdCommit},
			{"runc", info.RuncCommit},
		} {
			fmt.Fprintf(dockerCli.Out(), "%s version: %s", ci.Name, ci.Commit.ID)
			if ci.Commit.ID != ci.Commit.Expected {
				fmt.Fprintf(dockerCli.Out(), " (expected: %s)", ci.Commit.Expected)
			}
			fmt.Fprintf(dockerCli.Out(), "\n")
		}
		if len(info.SecurityOptions) != 0 {
			kvs, err := types.DecodeSecurityOptions(info.SecurityOptions)
			if err != nil {
				return err
			}
			fmt.Fprintf(dockerCli.Out(), "Security Options:\n")
			for _, so := range kvs {
				fmt.Fprintf(dockerCli.Out(), " %s\n", so.Name)
			}
		}
	}

	fprintfIfNotEmpty(dockerCli.Out(), "Kernel Version: %s\n", info.KernelVersion)
	fprintfIfNotEmpty(dockerCli.Out(), "Operating System: %s\n", info.OperatingSystem)
	fprintfIfNotEmpty(dockerCli.Out(), "OSType: %s\n", info.OSType)
	fprintfIfNotEmpty(dockerCli.Out(), "Architecture: %s\n", info.Architecture)
	fmt.Fprintf(dockerCli.Out(), "CPUs: %d\n", info.NCPU)
	fmt.Fprintf(dockerCli.Out(), "Total Memory: %s\n", units.BytesSize(float64(info.MemTotal)))
	fprintfIfNotEmpty(dockerCli.Out(), "Name: %s\n", info.Name)
	fprintfIfNotEmpty(dockerCli.Out(), "ID: %s\n", info.ID)
	fmt.Fprintf(dockerCli.Out(), "Docker Root Dir: %s\n", info.DockerRootDir)
	fmt.Fprintf(dockerCli.Out(), "Debug Mode (client): %v\n", debug.IsEnabled())
	fmt.Fprintf(dockerCli.Out(), "Debug Mode (server): %v\n", info.Debug)

	if info.Debug {
		fmt.Fprintf(dockerCli.Out(), " File Descriptors: %d\n", info.NFd)
		fmt.Fprintf(dockerCli.Out(), " Goroutines: %d\n", info.NGoroutines)
		fmt.Fprintf(dockerCli.Out(), " System Time: %s\n", info.SystemTime)
		fmt.Fprintf(dockerCli.Out(), " EventsListeners: %d\n", info.NEventsListener)
	}

	fprintfIfNotEmpty(dockerCli.Out(), "Http Proxy: %s\n", info.HTTPProxy)
	fprintfIfNotEmpty(dockerCli.Out(), "Https Proxy: %s\n", info.HTTPSProxy)
	fprintfIfNotEmpty(dockerCli.Out(), "No Proxy: %s\n", info.NoProxy)

	if info.IndexServerAddress != "" {
		u := dockerCli.ConfigFile().AuthConfigs[info.IndexServerAddress].Username
		if len(u) > 0 {
			fmt.Fprintf(dockerCli.Out(), "Username: %v\n", u)
		}
		fmt.Fprintf(dockerCli.Out(), "Registry: %v\n", info.IndexServerAddress)
	}

	if info.Labels != nil {
		fmt.Fprintln(dockerCli.Out(), "Labels:")
		for _, attribute := range info.Labels {
			fmt.Fprintf(dockerCli.Out(), " %s\n", attribute)
		}
		// TODO: Engine labels with duplicate keys has been deprecated in 1.13 and will be error out
		// after 3 release cycles (17.12). For now, a WARNING will be generated. The following will
		// be removed eventually.
		labelMap := map[string]string{}
		for _, label := range info.Labels {
			stringSlice := strings.SplitN(label, "=", 2)
			if len(stringSlice) > 1 {
				// If there is a conflict we will throw out a warning
				if v, ok := labelMap[stringSlice[0]]; ok && v != stringSlice[1] {
					fmt.Fprintln(dockerCli.Err(), "WARNING: labels with duplicate keys and conflicting values have been deprecated")
					break
				}
				labelMap[stringSlice[0]] = stringSlice[1]
			}
		}
	}

	fmt.Fprintf(dockerCli.Out(), "Experimental: %v\n", info.ExperimentalBuild)
	if info.ClusterStore != "" {
		fmt.Fprintf(dockerCli.Out(), "Cluster Store: %s\n", info.ClusterStore)
	}

	if info.ClusterAdvertise != "" {
		fmt.Fprintf(dockerCli.Out(), "Cluster Advertise: %s\n", info.ClusterAdvertise)
	}

	if info.RegistryConfig != nil && (len(info.RegistryConfig.InsecureRegistryCIDRs) > 0 || len(info.RegistryConfig.IndexConfigs) > 0) {
		fmt.Fprintln(dockerCli.Out(), "Insecure Registries:")
		for _, registry := range info.RegistryConfig.IndexConfigs {
			if !registry.Secure {
				fmt.Fprintf(dockerCli.Out(), " %s\n", registry.Name)
			}
		}

		for _, registry := range info.RegistryConfig.InsecureRegistryCIDRs {
			mask, _ := registry.Mask.Size()
			fmt.Fprintf(dockerCli.Out(), " %s/%d\n", registry.IP.String(), mask)
		}
	}

	if info.RegistryConfig != nil && len(info.RegistryConfig.Mirrors) > 0 {
		fmt.Fprintln(dockerCli.Out(), "Registry Mirrors:")
		for _, mirror := range info.RegistryConfig.Mirrors {
			fmt.Fprintf(dockerCli.Out(), " %s\n", mirror)
		}
	}

	printStorageDriverWarnings(dockerCli, info)

	return nil
}

func printStorageDriverWarnings(dockerCli *command.DockerCli, info types.Info) {
	if info.DriverStatus == nil {
		return
	}

	for _, pair := range info.DriverStatus {
		if pair[0] == "Data loop file" {
			fmt.Fprintf(dockerCli.Err(), "WARNING: %s: usage of loopback devices is "+
				"strongly discouraged for production use.\n         "+
				"Use `--storage-opt dm.thinpooldev` to specify a custom block storage device.\n", info.Driver)
		}
		if pair[0] == "Supports d_type" && pair[1] == "false" {
			backingFs := getBackingFs(info)

			msg := fmt.Sprintf("WARNING: %s: the backing %s filesystem is formatted without d_type support, which leads to incorrect behavior.\n", info.Driver, backingFs)
			if backingFs == "xfs" {
				msg += "         Reformat the filesystem with ftype=1 to enable d_type support.\n"
			}
			msg += "         Running without d_type support will not be supported in future releases."
			fmt.Fprintln(dockerCli.Err(), msg)
		}
	}
}

func getBackingFs(info types.Info) string {
	if info.DriverStatus == nil {
		return ""
	}

	for _, pair := range info.DriverStatus {
		if pair[0] == "Backing Filesystem" {
			return pair[1]
		}
	}
	return ""
}

func formatInfo(dockerCli *command.DockerCli, info types.Info, format string) error {
	tmpl, err := templates.Parse(format)
	if err != nil {
		return cli.StatusError{StatusCode: 64,
			Status: "Template parsing error: " + err.Error()}
	}
	err = tmpl.Execute(dockerCli.Out(), info)
	dockerCli.Out().Write([]byte{'\n'})
	return err
}

func fprintfIfNotEmpty(w io.Writer, format, value string) (int, error) {
	if value != "" {
		return fmt.Fprintf(w, format, value)
	}
	return 0, nil
}
