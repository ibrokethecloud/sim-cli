package cmd

import (
	"context"

	"github.com/ibrokethecloud/sim-cli/pkg/docker"
)

type Simulator struct {
	Name         string
	BundlePath   string
	Status       string
	Port         int
	Ctx          context.Context
	Image        string
	DockerClient docker.Client
}
