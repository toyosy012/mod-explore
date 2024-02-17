package usecase

import (
	"testing"

	"github.com/samber/do"
	"github.com/stretchr/testify/suite"

	"mods-explore/ark/omega/logic/variant/domain/service"
)

type VariantTestSuite struct {
	suite.Suite

	mockDB  *mockDBClient
	usecase VariantUsecase
}

func newTestVariantSuite() *VariantTestSuite { return &VariantTestSuite{} }

func TestVariantSuite(t *testing.T) {
	suite.Run(t, newTestVariantSuite())
}

func (s *VariantTestSuite) SetupSuite() {
	injector := do.New()

	mockDB := newMockDBClient()
	do.ProvideValue[service.VariantRepository](injector, mockDB)
	s.mockDB = mockDB
	usecase, err := NewVariant(injector)
	if err != nil {
		return
	}

	s.usecase = usecase
}
