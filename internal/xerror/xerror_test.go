package xerror_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/anditakaesar/uwa-go-fullstack/internal/xerror"
	"github.com/stretchr/testify/assert"
)

func TestDefineStatusCode(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{
			name: "session error",
			err:  &xerror.ErrorSession{Message: "session"},
			want: http.StatusUnauthorized,
		},
		{
			name: "permission error",
			err:  &xerror.ErrorPermission{Message: "permission"},
			want: http.StatusForbidden,
		},
		{
			name: "not found error",
			err:  &xerror.ErrorNotFound{Message: "not found"},
			want: http.StatusNotFound,
		},
		{
			name: "bad req error",
			err:  &xerror.ErrorBadRequest{Message: "bad req"},
			want: http.StatusBadRequest,
		},
		{
			name: "default error",
			err:  errors.New("some-error"),
			want: http.StatusInternalServerError,
		},
		{
			name: "no error",
			err:  nil,
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := xerror.DefineStatusCode(tt.err)
			assert.Equal(t, tt.want, got)
		})
	}
}
