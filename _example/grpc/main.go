package main

import (
	"github.com/Ja7ad/forker"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"log"
)

func main() {
	f := forker.NewGrpcForker(nil)
	srv := f.GetGrpcServer()

	grpc_health_v1.RegisterHealthServer(srv, health.NewServer())

	reflection.Register(srv)

	log.Fatalln(f.ServeGrpc(":9090"))
}
