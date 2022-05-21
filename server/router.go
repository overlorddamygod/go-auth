package server

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/overlorddamygod/go-auth/configs"
	"go.uber.org/fx"
)

func NewRouter(lc fx.Lifecycle, config *configs.Config) *gin.Engine {
	router := gin.New()

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			PORT := config.PORT

			if err := router.Run(":" + PORT); err != nil {
				return err
			}
			log.Println("Server started on port " + PORT)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("Stopping Auth server.")
			return nil
		},
	})
	return router
}
