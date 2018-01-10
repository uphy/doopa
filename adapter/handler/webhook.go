package handler

import (
	"fmt"
	"io/ioutil"

	"github.com/labstack/echo"
	"github.com/uphy/doopa/domain"
)

type (
	WebHookHandler struct {
		projectRepository domain.ProjectRepository
		deployService     domain.DeployService
	}
)

func NewWebHook(projectRepository domain.ProjectRepository, deployService domain.DeployService) *WebHookHandler {
	return &WebHookHandler{projectRepository, deployService}
}

func (w *WebHookHandler) Gogs(c echo.Context) error {
	body := c.Request().Body
	defer body.Close()
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	/*
		webhook := usecase.NewWebHook(w.projectRepository, w.deployService)

		err := webhook.WebHook(group, name, repoURL)
		if err != nil {
			return echo.NewHTTPError(500, err)
		}
	*/
	return nil
}
