package repositories

import (
	"Praetor/internal/models"
	"context"
	"log"
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
