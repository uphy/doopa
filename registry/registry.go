package registry

import (
	"github.com/uphy/doopa/domain"
	"github.com/uphy/doopa/infra/docker"
	"github.com/uphy/doopa/infra/git"
)

type (
	Registry struct {
		ProjectRepository domain.ProjectRepository
		DeployService     domain.DeployService
	}
)

func NewRegistry(dir string) (*Registry, error) {
	projectRepository := git.NewProjectRepository(dir)
	deployService, err := docker.NewDeployService()
	if err != nil {
		return nil, err
	}
	return &Registry{projectRepository, deployService}, nil
}
