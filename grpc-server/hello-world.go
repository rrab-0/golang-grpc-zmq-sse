package grpc_server

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"grpc-zmq-sse/db"
	pb "grpc-zmq-sse/generated-proto"
	zmq_local "grpc-zmq-sse/zmq-local"

	"google.golang.org/grpc"
)

var (
	portEnv, _ = strconv.Atoi(os.Getenv("GRPC_PORT"))
	port       = flag.Int("port", portEnv, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("GRPC Server received: %v", in.GetName())
	zmq_local.GlobalPublisher.Send("10001 "+in.GetName(), 1)
	log.Println("ZMQ PUB Sent: " + in.GetName())

	err := db.GlobalConnection.Create(&db.Dump{Message: in.GetName()}).Error
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	log.Println("PostgreSQL at grpc-server received: " + in.GetName())

	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func Start() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})

	log.Printf("grpc server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
