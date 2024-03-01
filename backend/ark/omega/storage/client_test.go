package storage

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do"
	"github.com/stretchr/testify/suite"

	"mods-explore/ark/omega"
)

type testModel struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type TestClientSuite struct {
	suite.Suite

	cli *Client
}

func TestTestClientSuite(t *testing.T) {
	suite.Run(t, &TestClientSuite{})
}

func (s *TestClientSuite) SetupSuite() {
	injector := do.New()
	{

		do.Provide(injector, func(_ *do.Injector) (omega.Environments, error) {
			conf, err := omega.LoadConfig()
			return *conf, err
		})
		do.ProvideValue(injector, slog.New(slog.NewJSONHandler(os.Stdout, nil)))

		do.Provide(injector, func(i *do.Injector) (*sqlx.DB, error) {
			env := do.MustInvoke[omega.Environments](i)
			postgresDSN := fmt.Sprintf(
				"postgres://%s:%s@%s:%d/%s?sslmode=disable",
				env.DBUsername,
				env.DBPassword,
				env.DatabaseURL,
				env.Port,
				env.DatabaseName,
			)
			return ConnectPostgres(postgresDSN)
		})

		do.Provide(injector, NewSQLxClient)
	}

	db := do.MustInvoke[*Client](injector)
	s.cli = db
	{ // migration up
		driver, err := postgres.WithInstance(db.DB.DB, &postgres.Config{})
		if err != nil {
			s.T().Log(fmt.Printf("error getting driver: %s", err.Error()))
			return
		}

		if err = RunMigration(driver, MigrateUp()); err != nil {
			if err != migrate.ErrNoChange {
				s.T().Log(fmt.Printf("error migrate up: %s", err.Error()))
				return
			}
		}

		if err = RunMigration(driver, VerifyMigrationVersion(migrationVer)); err != nil {
			s.T().Log(fmt.Printf("error migrate up: %s", err.Error()))
			return
		}
		s.T().Log("create table: test")
	}
}

func (s *TestClientSuite) TearDownSuite() {
	{ // migration down
		driver, err := postgres.WithInstance(s.cli.DB.DB, &postgres.Config{})
		if err != nil {
			s.T().Log(fmt.Printf("error getting driver: %s", err.Error()))
			return
		}
		if err = RunMigration(driver, MigrateDown()); err != nil {
			s.T().Log(fmt.Printf("error migrate down: %s", err.Error()))
			return
		}
		s.T().Log("delete table: test")
	}

	if err := s.cli.Close(); err != nil {
		s.T().Log(fmt.Printf("error db close: %s", err.Error()))
		return
	}
}

func (s *TestClientSuite) TestNamedGet() {
	ctx := context.Background()
	timeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	r, err := NamedGet[testModel](timeout, s.cli, `SELECT id, name FROM tests WHERE id = :id`, map[string]any{"id": 1})
	if err != nil {
		s.T().Log(fmt.Printf("failed select test record: %s", err.Error()))
		s.T().Fail()
	}

	s.Equal(&testModel{ID: 1, Name: "test"}, r)
}

func (s *TestClientSuite) TestSelect() {
	ctx := context.Background()
	timeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	rs, err := Select[testModel](timeout, s.cli, `SELECT id, name FROM tests`)
	if err != nil {
		s.T().Log(fmt.Printf("failed select test record: %s", err.Error()))
		s.T().Fail()
	}

	s.Equal([]testModel{{ID: 1, Name: "test"}}, rs)
}

func (s *TestClientSuite) TestNamedStore() {
	ctx := context.Background()
	timeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	id, err := NamedStore[int](timeout, s.cli, `INSERT INTO tests(name) VALUES(:name) RETURNING id`, map[string]any{"name": "test2"})
	if err != nil {
		s.T().Log(fmt.Printf("insert test record: %s", err.Error()))
		s.T().Fail()
	}

	r, err := NamedGet[testModel](timeout, s.cli, `SELECT id, name FROM tests WHERE id = :id`, map[string]any{"id": id})
	if err != nil {
		s.T().Log(fmt.Printf("failed select test record: %s", err.Error()))
		s.T().Fail()
	}

	s.Equal(&testModel{ID: id, Name: "test2"}, r)
}

func (s *TestClientSuite) TestNamedDelete() {
	ctx := context.Background()
	timeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	id, err := NamedStore[int](timeout, s.cli, `INSERT INTO tests(name) VALUES(:name) RETURNING id`, map[string]any{"name": "test3"})
	if err != nil {
		s.T().Log(fmt.Printf("insert test record: %s", err.Error()))
		s.T().Fail()
	}

	err = NamedDelete(timeout, s.cli, `DELETE FROM tests WHERE id = :id`, map[string]any{"id": id})
	s.ErrorIs(err, nil)
}
