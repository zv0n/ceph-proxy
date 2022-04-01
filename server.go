package main

import (
	"log"
	"net"
	"os"

	"github.com/zv0n/ceph-proxy/cephrpc"
	"google.golang.org/grpc"
)

func main() {
	// TODO configurable
	socketaddr := "/tmp/ceph-proxy.sock"
	listener, err := net.Listen("unix", socketaddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	if err := os.Remove(socketaddr); err != nil && !os.IsNotExist(err) {
		log.Fatalf("Failed to remove %s, error: %s", socketaddr, err.Error())
	}

	cephServer := cephrpc.Server{}

	grpcServer := grpc.NewServer()

	cephrpc.RegisterMountServiceServer(grpcServer, &cephServer)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
