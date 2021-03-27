package grpcerrwrap_test

import (
	"errors"
	"fmt"

	pkgerrors "github.com/pkg/errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"grpcerrwrap"
)

func ExampleCode() {
	err := loseStackTrace()
	fmt.Println(status.Convert(err).Code().String(), getStackHead(err))

	err = useGRPCWrapErr()
	fmt.Println(status.Convert(err).Code().String(), getStackHead(err))
	// Output:
	// InvalidArgument no stack trace
	// InvalidArgument example_test.go:46
}

func useGRPCWrapErr() error {
	err := someApplicationMethod()
	if err != nil {
		// we want to respond with a specific-error code, but not lose stack-traces for
		// our callers and interceptors. So we can use grpcerrwrap
		err = grpcerrwrap.Code(err, codes.InvalidArgument)
	}
	return err
}

func loseStackTrace() error {
	err := someApplicationMethod()
	if err != nil {
		err = status.New(codes.InvalidArgument, err.Error()).Err()
	}
	return err
}

func someApplicationMethod() error {
	// something goes wrong, we return an error with a useful stack-trace
	return pkgerrors.New("boom")
}

func getStackHead(err error) string {
	var hasStack interface {
		StackTrace() pkgerrors.StackTrace
	}
	if errors.As(err, &hasStack) {
		// return first frame if present
		for _, frame := range hasStack.StackTrace() {
			return fmt.Sprintf("%v", frame)
		}
		return "empty stack trace returned"
	}
	return "no stack trace"
}
