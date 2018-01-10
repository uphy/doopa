package docker

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types/filters"
	"github.com/pkg/errors"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/jhoonb/archivex"
	"github.com/uphy/doopa/domain"
)

type (
	deployService struct {
		ctx    context.Context
		docker *client.Client
	}
)

// NewDeployService creates DeployService.
func NewDeployService() (domain.DeployService, error) {
	defaultHeaders := map[string]string{"User-Agent": "doopa"}
	docker, err := client.NewClient("unix:///var/run/docker.sock", "", nil, defaultHeaders)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	return &deployService{ctx, docker}, nil
}

func (s *deployService) Deploy(id string, path string) error {
	// Check if already deployed or not.
	if err := s.assertNotDeployed(id); err != nil {
		return err
	}

	// Create build context archive
	tmpFile := filepath.Join(os.TempDir(), fmt.Sprintf("%d.tar", time.Now().Unix()))
	defer os.Remove(tmpFile)

	tar := new(archivex.TarFile)
	tar.Create(tmpFile)
	tar.AddAll(path, false)
	tar.Close()

	dockerBuildContext, err := os.Open(tmpFile)
	if err != nil {
		return err
	}
	defer dockerBuildContext.Close()

	// Build image
	resp, err := s.docker.ImageBuild(s.ctx, dockerBuildContext, types.ImageBuildOptions{
		SuppressOutput: false,
		Remove:         true,
		ForceRemove:    true,
		PullParent:     true,
		Tags:           []string{id},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if _, err := ioutil.ReadAll(resp.Body); err != nil {
		return err
	}

	// Delete the previous container
	if err := s.undeployIfExist(id); err != nil {
		return err
	}

	// Run the image
	if _, err := s.docker.ContainerCreate(s.ctx, &container.Config{
		Image: id,
	}, &container.HostConfig{}, &network.NetworkingConfig{}, id); err != nil {
		return err
	}
	if err := s.Start(id); err != nil {
		return err
	}
	return err
}

func (s *deployService) undeployIfExist(id string) error {
	if deployed, err := s.deployed(id); err != nil || !deployed {
		if err != nil {
			return err
		}
		return nil
	}
	return s.Undeploy(id)
}

func (s *deployService) Undeploy(id string) error {
	if err := s.assertDeployed(id); err != nil {
		return err
	}
	// remove container
	if err := s.docker.ContainerRemove(s.ctx, id, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	}); err != nil {
		return err
	}
	// remove image
	if _, err := s.docker.ImageRemove(s.ctx, id, types.ImageRemoveOptions{
		PruneChildren: true,
	}); err != nil {
		return errors.Wrap(err, "container successfully removed but image not")
	}
	return nil
}

func (s *deployService) assertNotDeployed(id string) error {
	if deployed, err := s.deployed(id); err != nil || deployed {
		if err != nil {
			return err
		}
		return errors.New("already deployed")
	}
	return nil
}

func (s *deployService) assertDeployed(id string) error {
	if deployed, err := s.deployed(id); err != nil || !deployed {
		if err != nil {
			return err
		}
		return errors.New("not deployed")
	}
	return nil
}

func (s *deployService) deployed(id string) (bool, error) {
	if status, err := s.Status(id); err != nil || status != domain.DeploymentStatusNotDeployed {
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func (s *deployService) Status(id string) (domain.DeploymentStatus, error) {
	args := filters.NewArgs()
	args.Add("name", id)
	containers, err := s.docker.ContainerList(s.ctx, types.ContainerListOptions{
		All:     true,
		Filters: args,
	})
	if err != nil {
		return domain.DeploymentStatusNotDeployed, err
	}
	switch len(containers) {
	case 0:
		return domain.DeploymentStatusNotDeployed, nil
	case 1:
		container := containers[0]
		return s.convertStatus(container.State), nil
	default:
		return domain.DeploymentStatusNotDeployed, errors.New("duplicated id")
	}
}

func (s *deployService) convertStatus(status string) domain.DeploymentStatus {
	switch status {
	case "running":
		return domain.DeploymentStatusRunning
	case "exited":
		return domain.DeploymentStatusStopped
	default:
		return domain.DeploymentStatusStopped
	}
}

func (s *deployService) Start(id string) error {
	if err := s.assertDeployed(id); err != nil {
		return err
	}
	return s.docker.ContainerStart(s.ctx, id, types.ContainerStartOptions{})
}

func (s *deployService) Stop(id string) error {
	if err := s.assertDeployed(id); err != nil {
		return err
	}
	return s.docker.ContainerStop(s.ctx, id, nil)
}
