package internal

import (
	"encoding/json"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestApiError_MarshalJSON(t *testing.T) {
	err := errors.New("test error")

	apiErr := NewApiError(404, "Not Found", err, "resource", "user", "id", "1234")
	result, _ := json.Marshal(apiErr)

	require.Equal(t, "{\"code\":404,\"status\":400,\"message\":\"Not Found\",\"error\":\"test error\",\"details\":{\"id\":\"1234\",\"resource\":\"user\"}}", string(result))
}

func TestValidateError_Assignable(t *testing.T) {
	var err = NewValidateError(errors.Errorf("validate failed"))

	var valErr *ValidateError
	require.True(t, errors.As(err, &valErr))
}

func TestNotFoundError_Assignable(t *testing.T) {
	var err = NewNotFoundError(errors.Errorf("not found"))

	var notFoundErr *NotFoundError
	require.True(t, errors.As(err, &notFoundErr))
}
