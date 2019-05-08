package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/siredmar/simplechat/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	text     string
	nickname string
	trigger  bool
}

func (s *server) Send(ctx context.Context, req *chatpb.SendRequest) (*empty.Empty, error) {
	fmt.Printf("Received request: %v\n", req)
	s.nickname = req.GetMsg().GetNickname()
	s.text = req.GetMsg().GetMessage()
	s.trigger = true
	return &empty.Empty{}, nil
}

// Streams back all other messages from the chat room
func (s *server) Receive(e *empty.Empty, stream chatpb.ChatService_ReceiveServer) error {
	for {
		if s.trigger == true {
			s.trigger = false
			err := stream.Send(&chatpb.ReceiveResponse{
				Msg: &chatpb.Message{
					Message:  s.text,
					Nickname: s.nickname,
				},
			})
			if err != nil {
				if err == io.EOF {
					// client has stopped receiving
					return nil
				}
				log.Fatalf("Fatal error: %v\n", err)
			}
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Error creating listener: %v\n", err)
	}
	s := grpc.NewServer()
	chatpb.RegisterChatServiceServer(s, &server{})

	// Register reflection service on gRPC server.
	reflection.Register(s)

	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("Error serving grpc server: %v\n:", err)
	}
}
