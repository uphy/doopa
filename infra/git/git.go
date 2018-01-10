package git

import (
	"io/ioutil"
	"os"
	"path/filepath"

	gg "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"

	"github.com/uphy/doopa/domain"
)

type (
	projectRepository struct {
		directory string
	}
)

// NewProjectRepository creates new ProjectRepository.
func NewProjectRepository(directory string) domain.ProjectRepository {
	return &projectRepository{directory}
}

func (s *projectRepository) New(group string, name string, remoteURL string) (*domain.Project, error) {
	dir := filepath.Join(s.directory, group, name)
	if err := os.MkdirAll(dir, 0777); err != nil {
		return nil, err
	}
	_, err := gg.PlainClone(dir, false, &gg.CloneOptions{
		URL:   remoteURL,
		Depth: 1,
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *projectRepository) Update(project *domain.Project) error {
	dir := filepath.Join(s.directory, project.Group, project.Name)
	repo, err := gg.PlainOpen(dir)
	if err != nil {
		return err
	}
	if err := repo.Fetch(&gg.FetchOptions{
		RemoteName: "origin",
		Depth:      1,
	}); err != nil && err != gg.NoErrAlreadyUpToDate {
		return err
	}
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}
	return worktree.Checkout(&gg.CheckoutOptions{
		Hash: plumbing.NewHash("origin/master"),
	})
}

func (s *projectRepository) Find(group string, name string) (*domain.Project, error) {
	dir := filepath.Join(s.directory, group, name)
	repo, err := gg.PlainOpen(dir)
	if err != nil {
		if err == gg.ErrRepositoryNotExists {
			return nil, nil
		}
		return nil, err
	}
	remote, err := repo.Remote("origin")
	url := remote.Config().URLs[0]
	return domain.NewProject(group, name, dir, domain.NewRemote(url)), nil
}

func (s *projectRepository) FindAll() ([]domain.Project, error) {
	groups, err := ioutil.ReadDir(s.directory)
	if err != nil {
		return nil, err
	}
	var projects []domain.Project
	for _, group := range groups {
		groupName := group.Name()
		groupPath := filepath.Join(s.directory, groupName)
		names, err := ioutil.ReadDir(groupPath)
		if err != nil {
			return nil, err
		}
		for _, name := range names {
			project, err := s.Find(groupName, name.Name())
			if err != nil {
				return nil, err
			}
			projects = append(projects, *project)
		}
	}
	return projects, nil
}

func (s *projectRepository) Delete(project *domain.Project) error {
	return os.RemoveAll(s.directory)
}
