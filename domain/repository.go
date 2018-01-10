package domain

type (
	// ProjectRepository is a repository interface for Project.
	ProjectRepository interface {
		// New clones the remote repository.
		New(group string, name string, remoteURL string) (*Project, error)
		// Update updates the local repository.
		Update(project *Project) error
		// Find finds a project by the group and name.
		Find(group string, name string) (*Project, error)
		// FindAll finds all projects.
		FindAll() ([]Project, error)
		// Delete deletes the project.
		Delete(project *Project) error
	}
)
