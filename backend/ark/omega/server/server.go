package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"

	"mods-explore/ark/omega"
	"mods-explore/ark/omega/logic/variant/usecase"
	"mods-explore/ark/omega/server/handlers"
	"mods-explore/ark/omega/storage"
)

func Run() {
	dbConf, err := omega.LoadConfig[omega.DBConfig]()
	if err != nil {
		logrus.Fatal(err)
		return
	}
	serverConf, err := omega.LoadConfig[omega.ServerConfig]()
	if err != nil {
		logrus.Fatal(err)
		return
	}
	s, err := newServer(*dbConf)
	if err != nil {
		logrus.Fatal(err)
		return
	}

	if err = s.Start(serverConf.Address); err != nil {
		logrus.Fatal(err)
	}
}

func newServer(conf omega.DBConfig) (*echo.Echo, error) {
	s := echo.New()
	s.HideBanner = true
	s.Use(middleware.Recover())
	s.Use(middleware.CORS())
	s.HTTPErrorHandler = handlers.ErrorHandler

	s.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "I'm fine!")
	})

	postgresDSN := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		conf.DBUsername,
		conf.DBPassword,
		conf.DatabaseURL,
		conf.Port,
		conf.DatabaseName,
	)

	variantsV1 := s.Group("/api/v1/variants")
	{ // variant
		repoClient, err := storage.NewVariantClient(postgresDSN)
		if err != nil {
			return nil, err
		}

		handler := handlers.NewVariant(usecase.NewVariant(repoClient))
		variantsV1.GET("/:id", handler.Read)
		variantsV1.GET("", handler.List)
		variantsV1.POST("/new", handler.Create)
		variantsV1.PUT("/:id", handler.Update)
		variantsV1.DELETE("/:id", handler.Delete)
	}

	return s, nil
}
