package testbed

import (
	"context"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun-starter-kit/bunapp"
	"github.com/uptrace/treemux"
)

func NewRequest(method, target string, body io.Reader) treemux.Request {
	return treemux.NewRequest(httptest.NewRequest(method, target, body))
}

func StartApp(t *testing.T) (context.Context, *bunapp.App) {
	ctx, app, err := bunapp.Start(context.TODO(), "test", "test")
	require.NoError(t, err)
	return ctx, app
}
