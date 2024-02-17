package usecase

import (
	"testing"

	"github.com/samber/do"
	"github.com/stretchr/testify/suite"

	"mods-explore/ark/omega/logic/creature/domain/service"
)

type UniqueDinosaurTestSuite struct {
	suite.Suite

	mockDB  *mockUniqueDB
	usecase UniqueUsecase
}

func newTestUniqueDinosaurSuite() *UniqueDinosaurTestSuite { return &UniqueDinosaurTestSuite{} }

func TestUniqueDinosaurSuite(t *testing.T) {
	suite.Run(t, newTestUniqueDinosaurSuite())
}

func (s *UniqueDinosaurTestSuite) SetupSuite() {
	injector := do.New()

	mockDB := newMockUniqueDB()
	do.ProvideValue[service.UniqueRepository](injector, mockDB)
	s.mockDB = mockDB
	usecase, err := NewUnique(injector)
	if err != nil {
		return
	}

	s.usecase = usecase
}
