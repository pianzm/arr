package main

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/pianzm/arr/config"
	"github.com/pianzm/arr/src/member/v1/delivery"
	usecase "github.com/pianzm/arr/src/member/v1/usecase"
)

func (s *HTTPServer) loadApp(app *echo.Echo, uc usecase.MemberUsecase) {
	memberDelivery := delivery.NewHandler(uc)
	memberGroup := app.Group("/v1/members/entity")
	memberDelivery.Mount(memberGroup)
}

func initHTTP(cfg *config.Config, uc usecase.MemberUsecase) {
	app := echo.New()
	app.Use(middleware.RequestID(), middleware.Recover(), middleware.Logger())

	httpServer := HTTPServer{
		uc: uc,
	}
	httpServer.loadApp(app, uc)

	app.Logger.Fatal(app.Start(fmt.Sprintf(":%d", cfg.DefaultPort)))
}
