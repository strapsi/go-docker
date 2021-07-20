package docker

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"

	_types "github.com/docker/docker/api/types"
	_container "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"

	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

type PsOptions struct {
	All         bool
	FilterNames []string
}

type RunOptions struct {
	Image string
	Name  string
	Force bool
	Env   map[string]string
	Ports map[string]string
}

var ctx context.Context
var cli *client.Client
var platform specs.Platform

func Ps(options *PsOptions) ([]_types.Container, error) {
	var filter filters.Args
	if len(options.FilterNames) > 0 {
		filter = filters.NewArgs()
		for _, f := range options.FilterNames {
			filter.Add("name", f)
		}
	}

	clo := _types.ContainerListOptions{}
	clo.Filters = filter
	clo.All = options.All

	containers, err := cli.ContainerList(ctx, clo)
	if err != nil {
		return nil, err
	}

	return containers, nil
}

func Run(options *RunOptions) error {
	ps_opts := PsOptions{}
	ps_opts.All = true
	ps_opts.FilterNames = append(ps_opts.FilterNames, options.Name)
	containers, err := Ps(&ps_opts)
	if err != nil {
		return err
	}

	var container_id string
	// container does not exist => create
	if len(containers) == 0 {
		env := make([]string, len(options.Env))
		i := 0
		for k, v := range options.Env {
			env[i] = fmt.Sprintf("%s=%s", k, v)
			i++
		}

		// create container
		config := _container.Config{
			Image: options.Image,
			Env:   env,
			ExposedPorts: map[nat.Port]struct{}{
				"9090/tcp": {},
			},
		}
		hostConfig := _container.HostConfig{
			PublishAllPorts: true,
			PortBindings: nat.PortMap{
				"9090/tcp": []nat.PortBinding{
					{
						HostPort: "9090",
						HostIP:   "0.0.0.0",
					},
				},
			},
		}
		networkConfig := network.NetworkingConfig{}
		result, err := cli.ContainerCreate(ctx, &config, &hostConfig, &networkConfig, &platform, options.Name)
		if err != nil {
			return err
		}
		container_id = result.ID
	} else { // container exists
		con := containers[0]
		if options.Force {
			container_id = con.ID
			switch con.State {
			case "running": // restart
				cli.ContainerStop(ctx, container_id, nil)
				// case "created":
				// case "exited":
				// default:
			}
		} else {
			return fmt.Errorf("container <%s> already exists: state = %s\nuse option `Force` or another container name\n", options.Name, con.State)
		}
	}

	// start container
	err = cli.ContainerStart(ctx, container_id, _types.ContainerStartOptions{})
	if err != nil {
		return err
	}
	return nil
}

func init() {
	var err error
	ctx = context.Background()
	cli, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	cli.NegotiateAPIVersion(ctx)

	platform = specs.Platform{}
	platform.OS = "linux"
}
