package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/jxohn/id-generate/consumer"
	pb "github.com/jxohn/id-generate/proto"
	"google.golang.org/grpc"
)

func main() {
	register, err := consumer.Register([]string{"localhost:2181"}, "/rpc/cn.bupt.john.test")
	if err != nil {
		log.Fatalf("failed to register consumer : %+v", err)
	}
	log.Println(register)

	for i := 0; i < 1000; i++ {
		intn := rand.Intn(len(register))
		dial, err := grpc.Dial(register[intn], grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("failed to connect : %+v", err)
		}

		defer dial.Close()

		client := pb.NewIDGenerateClient(dial)
		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
		defer cancelFunc()

		ds, err := client.GetIDs(ctx, &pb.GetIDsRequest{
			Biz: "FuckingBiz",
			Num: 100,
		})

		if err != nil {
			log.Fatalf("failed to get ids, %+v", err)
		}

		log.Printf("rsp is %+v", ds)
	}

}
