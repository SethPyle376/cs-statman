package grpcserver

import (
	"context"
	"net"
	"os"

	"github.com/sethpyle376/cs-statman/pkg/csproto"
	"github.com/sethpyle376/cs-statman/statman/data"
	"github.com/sethpyle376/cs-statman/statman/data/store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type GRPCServer struct {
	grpcs *grpc.Server
	csproto.UnimplementedStatmanServer
	errChan chan error
	db      data.Store
}

func New(port string) (*GRPCServer, error) {

	var grpcServer *grpc.Server

	cert, ok := os.LookupEnv("TLS_CERT_LOCATION")
	keyFile, keyOk := os.LookupEnv("TLS_KEY_LOCATION")

	if !ok || !keyOk {
		grpcServer = grpc.NewServer()
	} else {
		creds, err := credentials.NewServerTLSFromFile(cert, keyFile)
		if err != nil {
			return nil, err
		}
		println(creds.Info().ProtocolVersion)
		grpcServer = grpc.NewServer(grpc.Creds(creds))
	}

	gs := &GRPCServer{
		grpcs: grpcServer,
	}

	csproto.RegisterStatmanServer(gs.grpcs, gs)

	grpcPort, ok := os.LookupEnv("GRPC_PORT")

	if !ok {
		grpcPort = "4000"
	}

	lis, err := net.Listen("tcp", ":"+grpcPort)

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

func (gs *GRPCServer) GetPlayerMatchData(ctx context.Context, request *csproto.GetPlayerMatchDataRequest) (*csproto.GetPlayerMatchDataResponse, error) {
	response := &csproto.GetPlayerMatchDataResponse{}

	matchData, err := gs.db.GetPlayerMatchData(request.GetPlayerID())

	if err != nil {
		return nil, err
	}

	response.PlayerMatchData = matchData

	return response, nil
}

func (gs *GRPCServer) GetRecentMatches(ctx context.Context, request *csproto.GetRecentMatchesRequest) (*csproto.GetRecentMatchesResponse, error) {
	response := &csproto.GetRecentMatchesResponse{}

	matchData, err := gs.db.GetRecentMatches()

	if err != nil {
		return nil, err
	}

	response.MatchData = matchData

	return response, nil
}

func (gs *GRPCServer) GetPopularPlayers(ctx context.Context, request *csproto.GetPopularPlayersRequest) (*csproto.GetPopularPlayersResponse, error) {
	response := &csproto.GetPopularPlayersResponse{}

	playerData, err := gs.db.GetPopularPlayers()

	if err != nil {
		return nil, err
	}

	response.PopularPlayerData = playerData

	return response, nil
}
