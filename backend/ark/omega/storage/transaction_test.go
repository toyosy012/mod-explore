package storage

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"

	"mods-explore/ark/omega/logic/variant/domain/service"
)

type test struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type testTransactionSuite struct {
	suite.Suite

	cli  Client[test, int]
	mock sqlxmock.Sqlmock
}

func TestTransactionSuite(t *testing.T) {
	suite.Run(t, &testTransactionSuite{})
}

func (s *testTransactionSuite) SetupSuite() {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		s.T().Fatal(err)
	}

	s.cli = Client[test, int]{db, slog.New(slog.NewJSONHandler(os.Stdout, nil))}
	s.mock = mock
}

func (s *testTransactionSuite) committed(_ context.Context) (any, error) {
	s.T().Log("committed func")
	return &struct{}{}, nil
}
func (s *testTransactionSuite) panicked(_ context.Context) (any, error) {
	s.T().Log("panicked func")
	panic("test")
	return nil, service.IntervalServerError
}
func (s *testTransactionSuite) errored(_ context.Context) (any, error) {
	s.T().Log("errored func")
	return nil, service.NotFound
}

func (s *testTransactionSuite) TestCommitTransaction() {
	ctx := context.Background()

	s.mock.ExpectBegin()
	s.mock.ExpectCommit()

	r, _ := s.cli.WithTransaction(ctx, s.committed)
	s.Equal(&struct{}{}, r)
}

func (s *testTransactionSuite) TestPanicTransaction() {
	ctx := context.Background()

	s.mock.ExpectBegin()
	s.mock.ExpectRollback()

	_, err := s.cli.WithTransaction(ctx, s.panicked)
	s.ErrorIs(err, service.IntervalServerError)
}

func (s *testTransactionSuite) TestErrTransaction() {
	ctx := context.Background()

	s.mock.ExpectBegin()
	s.mock.ExpectRollback()

	_, err := s.cli.WithTransaction(ctx, s.errored)
	s.ErrorIs(err, service.NotFound)
}
