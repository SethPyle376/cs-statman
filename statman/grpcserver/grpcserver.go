package grpcserver

import (
	"context"
	"github.com/sethpyle376/cs-statman/pkg/csproto"
	"google.golang.org/grpc"
	"net"
)

type GRPCServer struct {
	grpcs *grpc.Server
	csproto.UnimplementedStatmanServer
	errChan chan error
}

func New(port string) (*GRPCServer, error) {
	gs := &GRPCServer{
		grpcs: grpc.NewServer(),
	}

	csproto.RegisterStatmanServer(gs.grpcs, gs)

	lis, err := net.Listen("tcp", ":4000")

	if err != nil {
		return nil, err
	}

	go func() {
		if err = gs.grpcs.Serve(lis); err != nil {
			gs.errChan <- err
		}
	}()
	return gs, nil
}

func (gs *GRPCServer) SaveMatch(ctx context.Context, request *csproto.SaveMatchRequest) (*csproto.SaveMatchResponse, error) {
	println("SAVING MATCH GOT VALUE: " + request.MatchInfo.MatchData.Map)
	return &csproto.SaveMatchResponse{}, nil
}

func (gs *GRPCServer) Stop() {}
