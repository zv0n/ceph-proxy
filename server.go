package main

import (
	"log"
	"net"

	"github.com/zv0n/ceph-proxy/cephrpc"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Faile to listen: %v", err)
	}

	s := cephrpc.Server{}

	grpcServer := grpc.NewServer()

	cephrpc.RegisterMountServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}

}
