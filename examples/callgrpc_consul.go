package main

import (
	"context"

	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/chaihaobo/be-template/proto"
)

func main() {
	conn, err := grpc.NewClient("consul://127.0.0.1:8500/betemplate-grpc?wait=14s",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := proto.NewHelloServiceClient(conn)
	response, err := client.SayHello(context.Background(), &proto.HelloRequest{
		Name: "Boice",
	})
	if err != nil {
		panic(err)
	}
	println(response.Reply)

}
