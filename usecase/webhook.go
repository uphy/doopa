package usecase

import (
	"fmt"

	"github.com/uphy/doopa/domain"
)

type (
	// WebHook represents the webhook usecase.
	WebHook struct {
		projectRepository domain.ProjectRepository
		deployService     domain.DeployService
	}
)

// NewWebHook creates new webhook usecase.
func NewWebHook(projectRepository domain.ProjectRepository, deployService domain.DeployService) *WebHook {
	return &WebHook{projectRepository, deployService}
}

// WebHook clone the git repository and deploy to the PaaS.
func (w *WebHook) WebHook(group, name, repoURL string) error {
	// find the project
	project, err := w.projectRepository.Find(group, name)
	if err != nil {
		return err
	}
	if project == nil {
		// if the project not exist, clone it.
		p, err := w.projectRepository.New(group, name, repoURL)
		if err != nil {
			return err
		}
		project = p
	} else {
		// if the project exists, update it.
		if err := w.projectRepository.Update(project); err != nil {
			return err
		}
	}

	// check the deployment status
	deployID := fmt.Sprintf("%s_%s", group, name)
	status, err := w.deployService.Status(deployID)
	if err != nil {
		return err
	}
	if status != domain.DeploymentStatusNotDeployed {
		// if the deployment has already been deployed, undeploy it first.
		err := w.deployService.Undeploy(deployID)
		if err != nil {
			return err
		}
	}
	// deploy
	return w.deployService.Deploy(deployID, project.Directory)
}
