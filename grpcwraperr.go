package grpcwraperr

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Code returns an error implementing GRPCStatus,
// which retains err. Will return nil if err is nil
// Status is formed from code and its message is err.Error()
func Code(err error, code codes.Code) error {
	if err == nil {
		return nil
	}
	return grpcStatusWrap{
		cause: err,
		status: status.New(code, err.Error()),
	}
}

// Status returns an error implementing GRPCStatus,
// which retains err. Will return nil if err is nil
// status should not be nil.
func Status(err error, status *status.Status) error {
	if err == nil {
		return nil
	}
	return grpcStatusWrap{
		cause:  err,
		status: status,
	}
}

type grpcStatusWrap struct {
	cause  error
	status *status.Status
}

func (g grpcStatusWrap) Error() string {
	return g.cause.Error()
}

var _ error = grpcStatusWrap{}

func (g grpcStatusWrap) Unwrap() error {
	return g.cause
}

func (g grpcStatusWrap) GRPCStatus() *status.Status {
	return g.status
}

