package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/internal/params"
)

func TestScenario(t *testing.T) {
	t.Run("should init scenario as started", func(t *testing.T) {
		assert.True(t, newScenario("test").HasStarted())
	})

	t.Run("should only create scenario if needed", func(t *testing.T) {
		store := NewScenarioStore()
		store.CreateNewIfNeeded("scenario-1")

		s, ok := store.FetchByName("scenario-1")
		assert.True(t, ok)
		assert.True(t, s.HasStarted())

		s.State = "another-state"
		store.Save(s)

		store.CreateNewIfNeeded("scenario-1")

		s, ok = store.FetchByName("scenario-1")
		assert.True(t, ok)
		assert.False(t, s.HasStarted())
		assert.Equal(t, s.State, "another-state")
	})
}

func TestScenarioConditions(t *testing.T) {
	store := NewScenarioStore()
	p := params.New()
	p.Set(BuiltInParamScenario, store)
	args := Args{Params: p}

	t.Run("should return true when scenario is not started and also not found", func(t *testing.T) {
		m := Scenario[any]("test", "required", "new")
		res, err := m(nil, args)

		assert.Nil(t, err)
		assert.True(t, res)
	})

	t.Run("should return false when scenario exists but it is not in the required state", func(t *testing.T) {
		store.CreateNewIfNeeded("hi")

		m := Scenario[any]("hi", "required", "new")
		res, err := m(nil, args)

		assert.Nil(t, err)
		assert.False(t, res)
	})
}
