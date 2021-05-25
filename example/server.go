package main

import (
	"context"
	"flag"
	"log"
	"net"

	pb "github.com/jxohn/rpc-register/proto"
	"github.com/jxohn/rpc-register/provider"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

type server struct {
	pb.UnimplementedIDGenerateServer
}

func (s *server) GetIDs(ctx context.Context, request *pb.GetIDsRequest) (*pb.GetIDsResponse, error) {
	log.Printf("receive request : %+v", request)

	return &pb.GetIDsResponse{
		Ids:     []int64{1234569084793287, 1234359873298698},
		ErrMsg:  "suc",
		ErrCode: 0,
	}, nil
}

func main() {
	var hostPort string
	flag.StringVar(&hostPort, "i", "", "服务ip端口")
	flag.Parse()
	log.Println(hostPort)
	err := provider.Online([]string{"localhost:2181"}, "/rpc/cn.bupt.john.test", hostPort)
	if err != nil {
		grpclog.Fatalf("failed to online %+v", err)
	}
	listen, err := net.Listen("tcp", hostPort)
	if err != nil {
		grpclog.Fatalf("failed to listen %+v", err)
	}

	s := grpc.NewServer()

	pb.RegisterIDGenerateServer(s, &server{})

	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve : %+v", err)
	}
}
