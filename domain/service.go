package domain

type (
	// DeploymentStatus is the status of the deployment.
	DeploymentStatus string

	// DeployService is a service interface for the deploment.
	DeployService interface {
		// Deploy deploys the deployment to the PaaS.
		Deploy(id string, path string) error
		// Undeploy removes deployment from the PaaS.
		Undeploy(id string) error
		// IsDeployed checks if the deployment exists in the PaaS.
		Status(id string) (DeploymentStatus, error)
		Start(id string) error
		Stop(id string) error
	}
)

const (
	// DeploymentStatusRunning represents the running status.
	DeploymentStatusRunning DeploymentStatus = "running"
	// DeploymentStatusStopped represents the stopped status.
	DeploymentStatusStopped DeploymentStatus = "stopped"
	// DeploymentStatusNotDeployed represents that the deployment is not the deployed.
	DeploymentStatusNotDeployed DeploymentStatus = "not deployed"
)
