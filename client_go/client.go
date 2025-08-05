package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/Anu-Ra-g/GRPC/coffeeshop_protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:9001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed to connect gRPC server")
	}

	defer conn.Close()

	c := pb.NewCoffeeShopClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	menuStream, err := c.GetMenu(ctx, &pb.MenuRequest{})
	if err != nil {
		log.Fatal("error calling function get menu")
	}

	done := make(chan bool)

	var items []*pb.Item

	go func() {
		for {
			resp, err := menuStream.Recv()
			if err == io.EOF {
				done <- true
				return
			}

			if err != nil {
				log.Fatalf("cannot recieve %v", err)
			}

			items = resp.Items
			log.Printf("Resp recieved: %v", items)
		}
	}()
	<-done

	receipt, err := c.PlaceOrder(ctx, &pb.Order{Items: items})
	log.Printf("%v", receipt)
	if err != nil {
		log.Fatalf("cannot recieve %v", err)
	}

	status, err := c.GetOrderStatus(ctx, receipt)
	log.Printf("%v", status)
	if err != nil {
		log.Fatalf("cannot recieve %v", err)
	}
}
