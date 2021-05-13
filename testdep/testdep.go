package testdep

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func Postgres(port int, pgUser string, pgPassword string, pgDatabase string) (stop func() error, err error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return
	}

	imageName := "postgres"

	_, err = cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return
	}

	natPort := nat.Port(fmt.Sprintf("%d/tcp", port))
	portSet := make(nat.PortSet)
	portSet[natPort] = struct{}{}
	portMap := make(nat.PortMap)
	portMap[natPort] = []nat.PortBinding{
		{
			HostIP:   "0.0.0.0",
			HostPort: strconv.Itoa(port),
		},
	}

	createConfig := &container.Config{
		Image:        imageName,
		ExposedPorts: portSet,
		Env:          []string{"POSTGRES_PASSWORD=12345678", "POSTGRES_USER=postgres", "POSTGRES_DB=postgres"},
	}
	hostConfig := &container.HostConfig{
		PortBindings: portMap,
	}

	resp, err := cli.ContainerCreate(ctx, createConfig, hostConfig, nil, nil, "")
	if err != nil {
		return
	}

	stop = func() error {
		td := 10 * time.Second
		err := cli.ContainerStop(ctx, resp.ID, &td)
		if err != nil {
			return err
		}
		err = cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		})
		if err != nil {
			return err
		}
		return nil
	}

	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		})
		return
	}

	return
}
