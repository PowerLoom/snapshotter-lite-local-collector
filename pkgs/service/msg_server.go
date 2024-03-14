package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"io"
	"net"
	"proto-snapshot-server/pkgs"
)

// server is used to implement submission.SubmissionService.
type server struct {
	pkgs.UnimplementedSubmissionServer
	stream network.Stream
}

var _ pkgs.SubmissionServer = &server{}

// NewMsgServerImpl returns an implementation of the SubmissionService interface
// for the provided Keeper.
func NewMsgServerImpl() pkgs.SubmissionServer {
	return &server{}
}

func setNewStream(s *server) error {
	st, err := rpctorelay.NewStream(network.WithUseTransient(context.Background(), "collect"), CollectorId, protocol.ConvertFromStrings([]string{"/collect"})[0])

	if err != nil {
		log.Debugln(err.Error())
		return errors.New("unable to establish stream")
	}
	s.stream = st

	return nil
}

func (s *server) SubmitSnapshot(stream pkgs.Submission_SubmitSnapshotServer) error {
	if s.stream == nil {
		err := setNewStream(s)
		log.Debugln(err)
	}
	for {
		submission, err := stream.Recv()

		if err == io.EOF {
			break
		} else if err != nil {
			log.Errorln("Grpc server crash ", err.Error())
			return err
		}

		log.Debugln("Received submission with request: ", submission.Request)

		subBytes, err := json.Marshal(submission)
		if err != nil {
			log.Debugln("Could not marshal submission")
			return err
		}
		if _, err = s.stream.Write(subBytes); err != nil {
			s.stream.Close()
			setNewStream(s)

			for i := 0; i < 5; i++ {
				_, err = s.stream.Write(subBytes)
				if err == nil {
					break
				} else {
					log.Errorln("Collector stream error, retrying: ", err.Error())
				}
			}
		}
	}
	return stream.SendAndClose(&pkgs.SubmissionResponse{Message: "Success"})
}

func (s *server) mustEmbedUnimplementedSubmissionServer() {
}

func StartSubmissionServer(server pkgs.SubmissionServer) {
	lis, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Debugf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pkgs.RegisterSubmissionServer(s, server)
	log.Debugln("Server listening at", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Debugf("failed to serve: %v", err)
	}
}