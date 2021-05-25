package main

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/jxohn/rpc-register/consumer"
	pb "github.com/jxohn/rpc-register/proto"
	"google.golang.org/grpc"
)

func main() {
	register, err := consumer.Register([]string{"localhost:2181"}, "/rpc/cn.bupt.john.test")
	if err != nil {
		log.Fatalf("failed to register consumer : %+v", err)
	}
	log.Println(register)

	for i := 0; i < 10000; i++ {
		conn, err := consumer.GetConn()
		if err != nil {
			log.Printf("failed to get conn, %+v\n", err)
			continue
		}
		dial, err := grpc.Dial(conn, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(500*time.Millisecond))
		if err != nil {
			log.Printf("failed to connect : %+v\n", err)
			continue
		}

		defer dial.Close()

		client := pb.NewIDGenerateClient(dial)
		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
		defer cancelFunc()

		ds, err := client.GetIDs(ctx, &pb.GetIDsRequest{
			Biz: "FuckingBiz_" + strconv.Itoa(i),
			Num: 100,
		})

		if err != nil {
			log.Fatalf("failed to get ids, %+v", err)
		}

		log.Printf("rsp is %+v", ds)
	}

}
