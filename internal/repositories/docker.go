package repositories

import (
	"Praetor/internal/models"
	"context"
	"io"
	"log"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type DockerRepository struct {
	Client *client.Client
	Ctx    context.Context
}

func NewDockerRepository(client *client.Client, ctx context.Context) *DockerRepository {
	return &DockerRepository{Client: client, Ctx: ctx}
}

func (d *DockerRepository) GetContainers() ([]models.Container, error) {
	containers, err := d.Client.ContainerList(d.Ctx, container.ListOptions{
		All: true,
	})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	parsedContainers := make([]models.Container, len(containers))
	for i, c := range containers {
		if c.State == "running" {
			parsedContainers[i], _ = d.GetContainer(c.ID)
		} else {
			parsedContainers[i] = summaryToLocalModel(c)
		}
	}
	return parsedContainers, nil
}

func (d *DockerRepository) GetContainer(id string) (models.Container, error) {
	cont, err := d.Client.ContainerInspect(d.Ctx, id)
	if err != nil {
		return models.Container{}, err
	}
	return inspectToLocalModel(cont), nil
}

func (d *DockerRepository) GetContainerLogs(id string, tail string) (string, error) {
	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: true,
		Tail:       tail,
	}

	logs, err := d.Client.ContainerLogs(d.Ctx, id, options)
	if err != nil {
		return "", err
	}
	defer logs.Close()

	buf := new(strings.Builder)
	_, err = io.Copy(buf, logs)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (d *DockerRepository) StopContainer(id string) error {
	return d.Client.ContainerStop(d.Ctx, id, container.StopOptions{})
}

func (d *DockerRepository) StartContainer(id string) error {
	return d.Client.ContainerStart(d.Ctx, id, container.StartOptions{})
}

func (d *DockerRepository) RestartContainer(id string) error {
	return d.Client.ContainerRestart(d.Ctx, id, container.StopOptions{})
}

func summaryToLocalModel(c container.Summary) models.Container {
	var ports []int
	for _, p := range c.Ports {
		if !slices.Contains(ports, int(p.PublicPort)) {
			ports = append(ports, int(p.PublicPort))
		}
	}
	return models.Container{
		ID:      c.ID,
		Image:   c.Image,
		Names:   c.Names,
		Ports:   ports,
		Created: time.Unix(c.Created, 0).Format(time.RFC3339),
		Status:  c.State,
	}
}

func inspectToLocalModel(cont container.InspectResponse) models.Container {
	var ports []int
	for portStr := range cont.Config.ExposedPorts {
		portNum := portStr.Port()
		if port, err := strconv.Atoi(portNum); err == nil {
			ports = append(ports, port)
		}
	}

	return models.Container{
		ID:      cont.ID,
		Image:   cont.Config.Image,
		Names:   []string{cont.Name},
		Ports:   ports,
		Created: cont.Created,
		Status:  cont.State.Status,
		Uptime:  getUptimeFromInspect(cont),
	}
}

func getUptimeFromInspect(inspect container.InspectResponse) string {
	if !inspect.State.Running {
		return ""
	}

	startedAt, err := time.Parse(time.RFC3339Nano, inspect.State.StartedAt)
	if err != nil {
		return ""
	}

	uptime := time.Since(startedAt)
	return uptime.Round(time.Second).String()
}
