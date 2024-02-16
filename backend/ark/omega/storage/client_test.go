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
	"github.com/stretchr/testify/suite"

	"mods-explore/ark/omega"
)

type testModel struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type TestClientSuite struct {
	suite.Suite

	cli *Client[testModel, int]
}

func TestTestClientSuite(t *testing.T) {
	suite.Run(t, &TestClientSuite{})
}

func (s *TestClientSuite) SetupSuite() {
	env, err := omega.LoadConfig()
	if err != nil {
		s.T().Log(fmt.Printf("error load db config: %s", err.Error()))
		return
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		env.DBUsername,
		env.DBPassword,
		env.DatabaseURL,
		env.Port,
		env.DatabaseName,
	)

	sqlxCli, err := ConnectPostgres(dsn)
	if err != nil {
		s.T().Log(fmt.Printf("error sqlx initialization: %s", err.Error()))
		return
	}

	{ // migration up
		driver, err := postgres.WithInstance(sqlxCli.DB, &postgres.Config{})
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

	s.cli, err = NewSQLxClient[testModel, int](
		sqlxCli,
		slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	)
	if err != nil {
		s.T().Log(fmt.Printf("error connection db: %s", err.Error()))
		return
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
	r, err := s.cli.NamedGet(timeout, `SELECT id, name FROM tests WHERE id = :id`, map[string]any{"id": 1})
	if err != nil {
		s.T().Log(fmt.Printf("failed select test record: %s", err.Error()))
		s.T().Fail()
	}

	s.Equal(&testModel{ID: 1, Name: "test"}, r)
}

func (s *TestClientSuite) TestNamedSelect() {
	ctx := context.Background()
	timeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	rs, err := s.cli.NamedSelect(timeout, `SELECT id, name FROM tests`)
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
	id, err := s.cli.NamedStore(timeout, `INSERT INTO tests(name) VALUES(:name) RETURNING id`, map[string]any{"name": "test2"})
	if err != nil {
		s.T().Log(fmt.Printf("insert test record: %s", err.Error()))
		s.T().Fail()
	}

	r, err := s.cli.NamedGet(timeout, `SELECT id, name FROM tests WHERE id = :id`, map[string]any{"id": id})
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
	id, err := s.cli.NamedStore(timeout, `INSERT INTO tests(name) VALUES(:name) RETURNING id`, map[string]any{"name": "test3"})
	if err != nil {
		s.T().Log(fmt.Printf("insert test record: %s", err.Error()))
		s.T().Fail()
	}

	err = s.cli.NamedDelete(timeout, `DELETE FROM tests WHERE id = :id`, map[string]any{"id": id})
	s.ErrorIs(err, nil)
}
