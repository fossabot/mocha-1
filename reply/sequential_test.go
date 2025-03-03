package reply

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/mock"
)

func TestSequential(t *testing.T) {
	t.Run("should return replies based configure sequence and return error when over", func(t *testing.T) {
		m := mock.New()
		m.Name = "mock_test"
		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
		builder := Sequential().
			Add(InternalServerError(), BadRequest(), OK(), NotFound())

		m.Hit()
		res, err := builder.Build(req, m, nil)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, res.Status)

		m.Hit()
		res, err = builder.Build(req, m, nil)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.Status)

		m.Hit()
		res, err = builder.Build(req, m, nil)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.Status)

		m.Hit()
		res, err = builder.Build(req, m, nil)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, res.Status)

		m.Hit()
		_, err = builder.Build(req, m, nil)
		assert.NotNil(t, err)
	})

	t.Run("should return replies based configure sequence and return error when over", func(t *testing.T) {
		m := mock.New()
		m.Name = "mock_test"
		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
		builder := Sequential().Add(OK()).ReplyOnSequenceEnded(NotFound())

		m.Hit()
		res, err := builder.Build(req, m, nil)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.Status)

		m.Hit()
		res, err = builder.Build(req, m, nil)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, res.Status)
	})
}

func TestShouldReturnErrorWhenSequenceDoesNotContainReplies(t *testing.T) {
	res, err := Sequential().Build(nil, nil, nil)
	assert.Nil(t, res)
	assert.NotNil(t, err)
}
