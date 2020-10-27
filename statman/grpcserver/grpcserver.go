package grpcserver

import (
	"context"
	"fmt"
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
	}

	println("db connection initialized")

	gs.db = db

	go func() {
		if err = gs.grpcs.Serve(lis); err != nil {
			gs.errChan <- err
		}
	}()
	return gs, nil
}

func (gs *GRPCServer) Stop() {}

func (gs *GRPCServer) SaveMatch(ctx context.Context, request *csproto.SaveMatchRequest) (*csproto.SaveMatchResponse, error) {
	error := gs.db.SaveMatch(request.MatchInfo)
	if error != nil {
		println(error.Error())
	}

	matchRequest := &csproto.GetMatchRequest{}
	matchRequest.MatchID = request.GetMatchInfo().GetMatchData().GetMatchID()

	response, err := gs.GetMatch(ctx, matchRequest)

	if err != nil {
		return nil, err
	}

	for _, player := range response.MatchInfo.PlayerData {
		fmt.Printf("NAME: %s\n", player.GetName())
		fmt.Printf("STEAMID: %d\n", player.GetSteamID())
		fmt.Printf("Kills: %d\n", player.GetKills())
		fmt.Printf("Deaths: %d\n", player.GetDeaths())
		fmt.Printf("ADR: %f\n\n\n", player.GetAdr())
	}

	return &csproto.SaveMatchResponse{}, error
}

func (gs *GRPCServer) GetPlayerMatchIDs(ctx context.Context, request *csproto.GetPlayerMatchIDsRequest) (*csproto.GetPlayerMatchIDsResponse, error) {
	matches, err := gs.db.GetPlayerMatches(request.GetPlayerID())

	if err != nil {
		return nil, err
	}

	response := &csproto.GetPlayerMatchIDsResponse{}
	response.MatchIDs = matches

	return response, nil
}

func (gs *GRPCServer) GetMatch(ctx context.Context, request *csproto.GetMatchRequest) (*csproto.GetMatchResponse, error) {
	match, err := gs.db.GetMatch(request.GetMatchID())
	if err != nil {
		return nil, err
	}

	response := &csproto.GetMatchResponse{}
	response.MatchInfo = match

	return response, nil
}
