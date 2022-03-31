package cephrpc

import (
	"log"

	"github.com/zv0n/ceph-proxy/ceph"
	"golang.org/x/net/context"
)

type Server struct {
}

func (s *Server) MountCeph(ctx context.Context, request *MountCephRequest) (*MountCephResponse, error) {
	log.Printf("Received request from client: client => \"%s\"; source => \"%s\"; target => \"%s\"; localUID => %d; remoteUID => %d; localGID => %d; remoteGID => %d", request.Client, request.MountSource, request.MountTarget, request.UidLocal, request.UidRemote, request.GidLocal, request.GidRemote)
	err := ceph.Mount(ceph.MountInput{
		Client:     request.Client,
		SourcePath: request.MountSource,
		TargetPath: request.MountTarget,
		UidLocal:   request.UidLocal,
		UidRemote:  request.UidRemote,
		GidLocal:   request.GidLocal,
		GidRemote:  request.GidRemote,
	})
	if err != nil {
		return &MountCephResponse{Output: "ERROR"}, err
	}
	return &MountCephResponse{Output: "Success"}, nil
}
