package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/siredmar/simplechat/proto"
	"google.golang.org/grpc"
)

func main() {
	waitc := make(chan bool)
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error connecting to grpc server: %v\n", err)
	}

	c := chatpb.NewChatServiceClient(cc)
	reader := bufio.NewReader(os.Stdin)

	go func() {
		for {
			text, _ := reader.ReadString('\n')
			// convert CRLF to LF
			text = strings.Replace(text, "\n", "", -1)

			_, err := c.Send(context.Background(), &chatpb.SendRequest{
				Msg: &chatpb.Message{
					Nickname: "User",
					Message:  text,
				},
			})
			if err != nil {
				log.Printf("Error: %v\n", err)
			}
		}
	}()

	go func() {
		stream, err := c.Receive(context.Background(), &empty.Empty{})
		if err != nil {
			fmt.Printf("Error making Receive call: %v\n", err)
		}
		for {
			res, err := stream.Recv()
			if err != nil {
				fmt.Printf("Error receiving stream: %v\n", err)
			}
			fmt.Println(res.GetMsg().GetNickname() + ": " + res.GetMsg().GetMessage())
		}
	}()
	<-waitc
}
