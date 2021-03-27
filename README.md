# grcpwraperr

gRPC wraps an original error while retaining the error-chain.

```go
wrapped := grpcerrwrap.Code(err, codes.InvalidArgument)
```

The returned `wrapped` above will perform the same in `errors.Is` and `errors.As` checks as `err`. As always,
we made sure to retain the chain. However, it will also return a status with code `codes.InvalidArgument` when used with the gRPC
status module, e.g `status.FromError(wrapped).Code()`. This means if we returned `wrapped` in a gRPC server method handler, the
calling client would receive `codes.InvalidArgument`:

```
func (s *SomeGRPCServer) SomeMethod(req *somepb.SomeRequest) (*somepb.SomeResponse, error) {
    // ...
    return grpcerrwrap.Code(err, codes.InvalidArgument)
}
```

The important thing, however, is that we could still retrieve the original error-chain of `err`. For instance we could
write an interceptor that sniffed for the `pkg/errors` stack-trace interface and report this in our error-tracker. 

## Full Example

Here is a [runnable example](example_test.go) of how `grpcerrwrap` retains the wrapped error's original error-chain, while working
with gRPC's methods.

```go
package grpcerrwrap_test

import (
	"errors"
	"fmt"

	pkgerrors "github.com/pkg/errors"
	"grpcerrwrap"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
```
