package cephrpc

import (
	"log"
	"os"

	"github.com/gogo/status"
	"github.com/zv0n/ceph-proxy/ceph"
	"golang.org/x/net/context"
	codes "google.golang.org/grpc/codes"
)

type Server struct {
}

func (s *Server) MountCeph(ctx context.Context, request *MountCephRequest) (*MountCephResponse, error) {
	log.Printf("Received mount request from client: client => \"%s\"; source => \"%s\"; target => \"%s\"; localUID => %d; remoteUID => %d; localGID => %d; remoteGID => %d", request.Client, request.MountSource, request.MountTarget, request.UidLocal, request.UidRemote, request.GidLocal, request.GidRemote)
	err, uidMap, gidMap := ceph.Mount(ceph.MountInput{
		Client:     request.Client,
		SourcePath: request.MountSource,
		TargetPath: request.MountTarget,
		UidLocal:   request.UidLocal,
		UidRemote:  request.UidRemote,
		GidLocal:   request.GidLocal,
		GidRemote:  request.GidRemote,
	})
	if err != nil {
		return &MountCephResponse{Output: err.Error(), UidMap: "", GidMap: ""}, err
	}
	return &MountCephResponse{Output: "Success", UidMap: uidMap, GidMap: gidMap}, nil
}

func (s *Server) UmountCeph(ctx context.Context, request *UmountCephRequest) (*UmountCephResponse, error) {
	log.Printf("Received umount request from client: target => %s", request.MountTarget)
	err := ceph.Umount(request.MountTarget)
	if err != nil {
		return &UmountCephResponse{Output: err.Error()}, err
	}
	err = os.Remove(request.UidMap)
	if err != nil {
		return &UmountCephResponse{Output: err.Error()}, status.Error(codes.Internal, err.Error())
	}
	err = os.Remove(request.GidMap)
	if err != nil {
		return &UmountCephResponse{Output: err.Error()}, status.Error(codes.Internal, err.Error())
	}
	return &UmountCephResponse{Output: "Success"}, nil
}
