package domain

type (
	// Remote is the information about the remote repository
	Remote struct {
		URL string
	}
	// Project is the project which is the work tree of Git repository.
	Project struct {
		// Group is the group name this project belongs to.
		Group string
		// Name is the repository name.
		Name string
		// Directory is the absolute directory path for this repository.
		Directory string
		// RemoteURL is the URL of the remote repository
		Remote *Remote
	}
	// Deployment is the deployment.
	Deployment struct {
		// ID is the unique string generated for each deployments.
		ID string
		// Path is the location of the deployment.
		Path string
	}
)

// NewRemote creates new remote repository information.
func NewRemote(url string) *Remote {
	return &Remote{url}
}

// NewProject create new local repository.
func NewProject(group string, name string, directory string, remote *Remote) *Project {
	return &Project{group, name, directory, remote}
}

// NewDeployment creates Deployment.
func NewDeployment(id string, path string) *Deployment {
	return &Deployment{id, path}
}
