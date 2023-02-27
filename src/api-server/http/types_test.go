package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Tests that checkQuery set the correct HTTP body.
func TestCheckQuery(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request, _ = http.NewRequest("GET", "http://example.com/?foo=bar&page=10&id=", nil)

	_, err := checkQuery(ctx, "non-existent-key")
	assert.Errorf(err, "key 'non-existent-key' does not exist")
	assert.Equal(`{"code":"500","message":"key 'non-existent-key' does not exist"}`, w.Body.String())

	val, err := checkQuery(ctx, "foo")
	assert.Nil(err)
	assert.Equal("bar", val)
}
