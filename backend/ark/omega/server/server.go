package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/samber/do"
	"github.com/sirupsen/logrus"

	"mods-explore/ark/omega"
	creatureUsecase "mods-explore/ark/omega/logic/creature/usecase"
	variantUsecase "mods-explore/ark/omega/logic/variant/usecase"
	"mods-explore/ark/omega/server/handlers"
	"mods-explore/ark/omega/storage"
)

func Run() {
	injector, err := Wired()
	if err != nil {
		logrus.Fatal(err)
		return
	}

	s, err := newServer(injector)
	if err != nil {
		logrus.Fatal(err)
		return
	}

	env := do.MustInvoke[omega.Environments](injector)
	if err = s.Start(env.Address); err != nil {
		logrus.Fatal(err)
	}
}

func newServer(injector *do.Injector) (*echo.Echo, error) {
	s := echo.New()
	s.HideBanner = true
	s.Use(middleware.Recover())
	s.Use(middleware.CORS())
	s.HTTPErrorHandler = handlers.NewErrorHandler(s)

	s.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "I'm fine!")
	})

	variantsV1 := s.Group(
		"/api/v1/variants",
		handlers.Transctioner(injector),
	)
	{ // variant
		handler := do.MustInvoke[handlers.VariantHandler](injector)
		variantsV1.GET("/:id", handler.Read)
		variantsV1.GET("", handler.List)
		variantsV1.POST("/new", handler.Create)
		variantsV1.PUT("/:id", handler.Update)
		variantsV1.DELETE("/:id", handler.Delete)
	}

	variantGroupsV1 := s.Group(
		"/api/v1/variant-groups",
		handlers.Transctioner(injector),
	)
	{ // variant group
		handler := do.MustInvoke[handlers.VariantGroupHandler](injector)
		variantGroupsV1.GET("/:id", handler.Read)
		variantGroupsV1.GET("", handler.List)
		variantGroupsV1.POST("/new", handler.Create)
		variantGroupsV1.PUT("/:id", handler.Update)
		variantGroupsV1.DELETE("/:id", handler.Delete)
	}
	{
		uniquesV1 := s.Group(
			"/api/v1/uniques",
			handlers.Transctioner(injector),
		)
		handler := do.MustInvoke[handlers.UniqueHandler](injector)
		uniquesV1.GET("/:id", handler.ReadUnique)
		uniquesV1.GET("", handler.ListUniques)
		uniquesV1.POST("/new", handler.CreateUnique)
		uniquesV1.PUT("/:id", handler.UpdateUnique)
		uniquesV1.DELETE("/:id", handler.DeleteUnique)
	}

	return s, nil
}

func Wired() (*do.Injector, error) {
	injector := do.New()

	do.Provide(injector, func(_ *do.Injector) (omega.Environments, error) {
		conf, err := omega.LoadConfig()
		return *conf, err
	})

	do.Provide(injector, func(i *do.Injector) (*sqlx.DB, error) {
		env := do.MustInvoke[omega.Environments](i)
		dns := fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=disable",
			env.DBUsername,
			env.DBPassword,
			env.DatabaseURL,
			env.Port,
			env.DatabaseName,
		)
		return storage.ConnectPostgres(dns)
	})

	do.ProvideValue(injector, slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	do.Provide(injector, storage.NewSQLxClient)

	do.Provide(injector, storage.NewVariantClient)
	do.Provide(injector, variantUsecase.NewVariant)
	do.Provide(injector, handlers.NewVariant)

	do.Provide(injector, storage.NewVariantGroupClient)
	do.Provide(injector, variantUsecase.NewVariantGroup)
	do.Provide(injector, handlers.NewVariantGroup)

	do.Provide(injector, storage.NewUniqueQueryRepo)
	do.Provide(injector, storage.NewUniqueCommandRepo)
	do.Provide(injector, storage.NewUniqueVariantsClient)
	do.Provide(injector, storage.NewDinosaurClient)
	do.Provide(injector, creatureUsecase.NewUnique)
	do.Provide(injector, handlers.NewUnique)

	return injector, nil
}
