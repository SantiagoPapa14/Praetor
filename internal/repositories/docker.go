package repositories

import (
	"Praetor/internal/models"
	"context"
	"log"
	"strconv"
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
		ports := make([]int, len(c.Ports))
		for j, p := range c.Ports {
			ports[j] = int(p.PublicPort)
		}

		parsedContainers[i] = models.Container{
			ID:      c.ID,
			Image:   c.Image,
			Names:   c.Names,
			Ports:   ports,
			Created: time.Unix(c.Created, 0).Format(time.RFC3339),
			Status:  c.State,
		}
	}

	return parsedContainers, nil
}

func (d *DockerRepository) GetContainer(id string) (models.Container, error) {
	cont, err := d.Client.ContainerInspect(d.Ctx, id)
	if err != nil {
		return models.Container{}, err
	}

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
	}, nil
}

func (d *DockerRepository) StopContainer(id string) error {
	return d.Client.ContainerStop(d.Ctx, id, container.StopOptions{})
}

func (d *DockerRepository) StartContainer(id string) error {
	return d.Client.ContainerStart(d.Ctx, id, container.StartOptions{})
}
