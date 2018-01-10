package docker

import (
	"testing"

	"github.com/uphy/doopa/domain"
)

func TestNewDeployService(t *testing.T) {
	_, err := NewDeployService()
	if err != nil {
		t.Error(err)
	}
}

func TestDeployServiceDeploy(t *testing.T) {
	service, err := NewDeployService()

	// deploy
	if err := service.Deploy("test", "/Users/ishikura/go/src/github.com/uphy/doopa/test/gitrepo"); err != nil {
		t.Fatal(err)
	}
	// test the status
	status, err := service.Status("test")
	if err != nil {
		t.Error(err)
	}
	if status != domain.DeploymentStatusRunning {
		t.Errorf("expected running status but %s.", status)
	}
	// undeploy
	if err := service.Undeploy("test"); err != nil {
		t.Error(err)
	}
}
