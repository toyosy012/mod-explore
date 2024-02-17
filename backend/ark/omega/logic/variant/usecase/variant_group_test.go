package usecase

import (
	"testing"

	"github.com/samber/do"
	"github.com/stretchr/testify/suite"

	"mods-explore/ark/omega/logic/variant/domain/service"
)

type VariantGroupTestSuite struct {
	suite.Suite

	mockDB  *mockVariantGroup
	usecase VariantGroupUsecase
}

func newTestVariantGroupSuite() *VariantGroupTestSuite { return &VariantGroupTestSuite{} }

func TestVariantGroupSuite(t *testing.T) {
	suite.Run(t, newTestVariantGroupSuite())
}

func (s *VariantGroupTestSuite) SetupSuite() {
	injector := do.New()

	mockDB := newMockVariantGroup()
	do.ProvideValue[service.VariantGroupRepository](injector, mockDB)
	s.mockDB = mockDB
	usecase, err := NewVariantGroup(injector)
	if err != nil {
		return
	}

	s.usecase = usecase
}
