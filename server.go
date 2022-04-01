package main

import (
	"log"
	"net"
	"os"

	"github.com/zv0n/ceph-proxy/cephrpc"
	"github.com/zv0n/ceph-proxy/configuration"
	"google.golang.org/grpc"
)

const configPath = "/etc/ceph-proxy.conf"

func main() {
	config, err := configuration.ParseConfigFile(configPath)
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("Could not read config file: \"%s\" - %v", configPath, err)
	}

	if err := os.Remove(config.SocketPath); err != nil && !os.IsNotExist(err) {
		log.Fatalf("Failed to remove %s, error: %s", config.SocketPath, err.Error())
	}

	listener, err := net.Listen("unix", config.SocketPath)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	cephServer := cephrpc.Server{
		Config: config,
	}

	grpcServer := grpc.NewServer()

	cephrpc.RegisterMountServiceServer(grpcServer, &cephServer)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
