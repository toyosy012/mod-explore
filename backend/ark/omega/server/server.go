package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/samber/do"
	"github.com/sirupsen/logrus"

	"mods-explore/ark/omega"
	"mods-explore/ark/omega/logic/variant/domain/service"
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
	s.HTTPErrorHandler = handlers.NewErrorHandler(s)

	s.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "I'm fine!")
	})

	postgresDSN := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		conf.DBUsername,
		conf.DBPassword,
		conf.DatabaseURL,
		conf.Port,
		conf.DatabaseName,
	)

	injector, err := Wired(postgresDSN)
	if err != nil {
		return nil, err
	}

	variantsV1 := s.Group("/api/v1/variants")
	{ // variant
		variantsV1.Use(
			handlers.Transctioner(
				do.MustInvoke[service.VariantRepository](injector).(storage.VariantClient).Client,
			),
		)
		variant := do.MustInvoke[usecase.VariantUsecase](injector)
		handler := handlers.NewVariant(variant)
		variantsV1.GET("/:id", handler.Read)
		variantsV1.GET("", handler.List)
		variantsV1.POST("/new", handler.Create)
		variantsV1.PUT("/:id", handler.Update)
		variantsV1.DELETE("/:id", handler.Delete)
	}

	variantGroupsV1 := s.Group("/api/v1/variant-groups")
	{ // variant group
		variantsV1.Use(
			handlers.Transctioner(
				do.MustInvoke[service.VariantGroupRepository](injector).(storage.VariantGroupClient).Client,
			),
		)
		variantGroup := do.MustInvoke[usecase.VariantGroupUsecase](injector)
		handler := handlers.NewVariantGroup(variantGroup)
		variantGroupsV1.GET("/:id", handler.Read)
		variantGroupsV1.GET("", handler.List)
		variantGroupsV1.POST("/new", handler.Create)
		variantGroupsV1.PUT("/:id", handler.Update)
		variantGroupsV1.DELETE("/:id", handler.Delete)
	}

	return s, nil
}

func Wired(postgresDSN string) (*do.Injector, error) {
	injector := do.New()

	db, err := storage.ConnectPostgres(postgresDSN)
	if err != nil {
		return nil, err
	}

	variantRepo, err := storage.NewVariantClient(db, slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	if err != nil {
		return nil, err
	}
	do.ProvideValue(injector, variantRepo)
	do.Provide(injector, usecase.NewVariant)

	variantGroupRepo, err := storage.NewVariantGroupClient(db, slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	if err != nil {
		return nil, err
	}
	do.ProvideValue(injector, variantGroupRepo)
	do.Provide(injector, usecase.NewVariantGroup)

	return injector, nil
}
