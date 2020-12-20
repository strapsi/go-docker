package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"

	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

type PsOptions struct {
	All bool
}

type RunOptions struct {
	Image string
	Name string
}

var ctx context.Context
var cli *client.Client
var platform specs.Platform

func Ps(options *PsOptions) ([]types.Container, error) {
	fmt.Printf("%s\n", options)
	clo := types.ContainerListOptions{}
	clo.All = options.All
	containers, err := cli.ContainerList(ctx, clo)
	if err != nil {
		return nil, err
	}

	if len(containers) > 0 {
		return containers, nil
	} else {
		return nil, nil
	}	
}

func Run(options *RunOptions) error {
	// for now: stop container
	// create container

	config := container.Config{}
	config.Image = options.Image
	hostConfig := container.HostConfig{}
	networkConfig := network.NetworkingConfig{}
	container, err := cli.ContainerCreate(ctx, &config, &hostConfig, &networkConfig, &platform, options.Name)
	if err != nil {
		return err
	}
	// start container

	err = cli.ContainerStart(ctx, container.ID, types.ContainerStartOptions{})
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

