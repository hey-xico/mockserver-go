package mockserver

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateSimpleRequest(t *testing.T) {

	req := ARequest().
		WithMethod("GET").
		WithHeader(map[string]string{
			"Content-Type": "application/json"}).
		WithPath("/foo")

	assert.Nil(t, req.body)
	assert.Nil(t, req.pathParameters)
	assert.Nil(t, req.queryParams)

	assert.Equal(t, "GET", req.method)
	assert.Equal(t, "/foo", req.path)
	assert.Equal(t, map[string]string{"Content-Type": "application/json"}, req.header)
}
