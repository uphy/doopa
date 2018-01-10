package git

import (
	"testing"
)

func TestProjectRepositoryNew(t *testing.T) {
	projectRepository := NewProjectRepository("repo")
	project, err := projectRepository.New("testgroup", "testproject", "../../test/gitrepo/")
	if err != nil {
		t.Error(err)
	}
	project, err = projectRepository.Find("testgroup", "testproject")
	if err := projectRepository.Update(project); err != nil {
		t.Error(err)
	}
	projects, err := projectRepository.FindAll()
	if err != nil {
		t.Error(err)
	}
	if len(projects) != 1 {
		t.Error("can't load project properly.")
	}
	if err := projectRepository.Delete(project); err != nil {
		t.Error(err)
	}
}
