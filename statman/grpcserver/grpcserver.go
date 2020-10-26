package grpcserver

import (
	"context"
	"github.com/sethpyle376/cs-statman/pkg/csproto"
	"github.com/sethpyle376/cs-statman/statman/data"
	"github.com/sethpyle376/cs-statman/statman/data/store"
	"google.golang.org/grpc"
	"net"
)

type GRPCServer struct {
	grpcs *grpc.Server
	csproto.UnimplementedStatmanServer
	errChan chan error
	db      data.Store
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

	db, err := store.New("postgres")

	if err != nil {
		return nil, err
	} else {
		println("db connection initialized")
	}

	gs.db = db

	go func() {
		if err = gs.grpcs.Serve(lis); err != nil {
			gs.errChan <- err
		}
	}()
	return gs, nil
}

func (gs *GRPCServer) SaveMatch(ctx context.Context, request *csproto.SaveMatchRequest) (*csproto.SaveMatchResponse, error) {
	error := gs.db.SaveMatch(request.MatchInfo)
	return &csproto.SaveMatchResponse{}, error
}

func (gs *GRPCServer) Stop() {}
