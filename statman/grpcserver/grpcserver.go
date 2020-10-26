package grpcserver

import (
	"context"
	"github.com/sethpyle376/cs-statman/pkg/csproto"
	"github.com/sethpyle376/cs-statman/statman/data"
	"github.com/sethpyle376/cs-statman/statman/data/store"
	"google.golang.org/grpc"
	"net"
	"strconv"
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
	for _, element := range request.MatchInfo.PlayerData {

		println(element.Name)
		println(strconv.Itoa(int(element.Team)))
		println("Kills: " + strconv.Itoa(int(element.Kills)))
		println("Deaths: " + strconv.Itoa(int(element.Deaths)))
		println("ADR: " + strconv.FormatFloat(float64(element.Adr), 'f', 2, 32) + "\n\n\n")
	}

	for index, round := range request.MatchInfo.RoundData {
		println("ROUND: " + strconv.Itoa(index))
		for _, kill := range round.Kills {
			println(strconv.FormatInt(int64(kill.KillerID), 10) + " killed " + strconv.FormatInt(int64(kill.VictimID), 10))
		}
		println("ROUND WON BY: " + strconv.Itoa(int(round.WinningTeam)))
	}
	return &csproto.SaveMatchResponse{}, nil
}

func (gs *GRPCServer) Stop() {}
