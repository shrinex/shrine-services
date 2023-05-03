package dtmx

import (
	"context"
	"github.com/dtm-labs/dtm/client/dtmcli"
	"github.com/dtm-labs/dtm/client/dtmgrpc"
)

func MustBarrierFromGrpc(ctx context.Context) *dtmcli.BranchBarrier {
	barrier, err := dtmgrpc.BarrierFromGrpc(ctx)
	if err != nil {
		panic(err)
	}

	return barrier
}
