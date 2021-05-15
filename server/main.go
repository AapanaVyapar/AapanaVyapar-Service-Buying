package main

import (
	"aapanavyapar-service-buying/pb"
	"aapanavyapar-service-buying/services"
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {
	log.Printf("Stating server on port  :  %d", os.Getenv("Port"))

	fmt.Println("Environmental Variables Loaded .. !!")

	service := services.NewBuyingService()

	grpcServer := grpc.NewServer()
	pb.RegisterBuyingServiceServer(grpcServer, service)

	address := fmt.Sprintf("0.0.0.0:%s", os.Getenv("Port"))
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Can not start server", err)
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("Can not start server", err)
	}

}
