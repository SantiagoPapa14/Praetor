package app

import "Praetor/internal/repositories"

type App struct {
	Repos struct {
		Session *repositories.SessionRepository
		User    *repositories.UserRepository
		Docker  *repositories.DockerRepository
	}
}
