package dtmx

import (
	"context"
	"fmt"
	"github.com/dtm-labs/dtm/client/dtmcli"
	"github.com/dtm-labs/dtm/client/dtmgrpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

func MustBarrierFromGrpc(ctx context.Context) *dtmcli.BranchBarrier {
	barrier, err := dtmgrpc.BarrierFromGrpc(ctx)
	if err != nil {
		panic(err)
	}

	return barrier
}

func MustBarrierFromQuery(r *http.Request) *dtmcli.BranchBarrier {
	barrier, err := dtmcli.BarrierFromQuery(r.URL.Query())
	if err != nil {
		panic(err)
	}

	return barrier
}

func Abort(err error) error {
	return asDtmError(err, codes.Aborted)
}

func Abortf(format string, a ...any) error {
	return status.Error(codes.Aborted, fmt.Sprintf(format, a))
}

func Ongoing(err error) error {
	return asDtmError(err, codes.FailedPrecondition)
}

func Ongoingf(format string, a ...any) error {
	return status.Error(codes.FailedPrecondition, fmt.Sprintf(format, a))
}

func Retry(err error) error {
	return asDtmError(err, codes.Internal)
}

func Retryf(format string, a ...any) error {
	return status.Error(codes.Internal, fmt.Sprintf(format, a))
}

func asDtmError(err error, code codes.Code) error {
	if s, ok := status.FromError(err); ok {
		if s.Code() == code {
			return err
		}

		return status.Error(code, s.Message())
	}

	return status.Error(code, err.Error())
}
