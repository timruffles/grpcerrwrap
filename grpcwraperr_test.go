package grpcerrwrap

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCompatibleWithStatusFromError(t *testing.T) {
	wrapped := Code(originalErr, codes.InvalidArgument)
	wrappedSt := Status(originalErr, status.New(codes.InvalidArgument, "status-msg"))

	assertCompatibleWithFromError := func(err error, expectedCode codes.Code) {
		st, ok := status.FromError(wrapped)
		if assert.True(t, ok) && assert.NotNil(t, st) {
			assert.Equal(t, expectedCode, st.Code())
		}
	}

	t.Run("Code", func(t *testing.T) {
		assertCompatibleWithFromError(wrapped, codes.InvalidArgument)
	})

	t.Run("Status", func(t *testing.T) {
		assertCompatibleWithFromError(wrappedSt, codes.InvalidArgument)
		assert.Equal(t, "original", wrappedSt.Error())
		assert.Equal(t, "status-msg", status.Convert(wrappedSt).Message())
	})
}

func TestUnwrappingWorks(t *testing.T) {
	wrapped := Code(originalErr, codes.InvalidArgument)
	assert.True(t, errors.Is(wrapped, originalErr))
}

func TestReturnsNilIfErrIsNil(t *testing.T) {
	assert.Nil(t, Code(nil, codes.InvalidArgument))
	assert.Nil(t, Status(nil, status.New(codes.InvalidArgument, "status-msg")))
}

var originalErr = errors.New("original")

type exampleStack struct {
}

