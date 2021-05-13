package testdep

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/jackc/pgconn"
)

var (
	HealthCheckTimeoutErr = errors.New("container timeout error")
)

type Postgres struct {
	Port     int
	User     string
	Password string
	Database string
}

func (c *Postgres) ConnString() string {
	return fmt.Sprintf("host=localhost port=%d user=%s password=%s dbname=%s", c.Port, c.User, c.Password, c.Database)
}

func (c *Postgres) Run(timeoutSeconds int) (stop func() error, err error) {
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

	portSet := nat.PortSet{"5432/tcp": struct{}{}}
	portMap := nat.PortMap{
		"5432/tcp": []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: strconv.Itoa(c.Port),
			},
		},
	}

	env := make([]string, 3)
	env[0] = fmt.Sprintf("POSTGRES_USER=%s", c.User)
	env[1] = fmt.Sprintf("POSTGRES_PASSWORD=%s", c.Password)
	env[2] = fmt.Sprintf("POSTGRES_DB=%s", c.Database)
	createConfig := &container.Config{
		Image:        imageName,
		ExposedPorts: portSet,
		Env:          env,
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

	wait := make(chan struct{})
	ok := false
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for range ticker.C {
			err := c.Ping()
			if err == nil {
				ok = true
				wait <- struct{}{}
				return
			}
			timeoutSeconds--
			if timeoutSeconds <= 0 {
				break
			}
		}
		wait <- struct{}{}
	}()
	<-wait
	if !ok {
		err = HealthCheckTimeoutErr
	}

	return
}

func (c *Postgres) Ping() error {
	pgConn, err := pgconn.Connect(context.Background(), c.ConnString())
	if err != nil {
		return fmt.Errorf("pgconn failed to connect: %s", err)
	}
	defer pgConn.Close(context.Background())
	result := pgConn.Exec(context.Background(), "SELECT 1")
	results, err := result.ReadAll()
	if err != nil {
		return fmt.Errorf("failed reading result: %s", err)
	}
	if len(results) != 1 {
		return fmt.Errorf("wrong result: %+v", results[0])
	}
	err = result.Close()
	if err != nil {
		return fmt.Errorf("failed reading result: %s", err)
	}
	return nil
}
