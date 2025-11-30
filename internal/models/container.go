package models

type Container struct {
	ID      string
	Image   string
	Names   []string
	Ports   []int
	Created string
	Status  string
	Uptime  string
}
