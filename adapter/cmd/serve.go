package cmd

import (
	"log"

	"github.com/labstack/echo"
	"github.com/spf13/cobra"
	"github.com/uphy/doopa/adapter/handler"
	"github.com/uphy/doopa/registry"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use: "serve",
		RunE: func(cmd *cobra.Command, args []string) error {
			e := echo.New()
			r, err := registry.NewRegistry("repo")
			if err != nil {
				panic(err)
			}

			webhook := handler.NewWebHook(r.ProjectRepository, r.DeployService)
			e.POST("/api/webhook/gogs", webhook.Gogs)

			if err := e.Start(":8080"); err != nil {
				log.Fatal(err)
			}
			return nil
		},
	})
}
