package grpc_server

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	zmq "github.com/pebbe/zmq4"

	"grpc-zmq-sse/db"
	pb "grpc-zmq-sse/generated-proto"
	zmq_local "grpc-zmq-sse/zmq-local"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("GRPC Server received: %v", in.GetName())
	zmq_local.GlobalPublisher.Send("10001 "+in.GetName(), zmq.DONTWAIT)
	log.Println("ZMQ PUB Sent: " + in.GetName())

	err := db.GlobalConnection.Create(&db.Dump{Message: in.GetName()}).Error
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	log.Println("PostgreSQL at grpc-server received: " + in.GetName())

	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func Start() {
	// flag.Parse()
	var (
		portEnvHello, _ = strconv.Atoi(os.Getenv("GRPC_PORT"))
		portHello       = flag.Int("portHello", portEnvHello, "The server port")
	)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *portHello))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	reflection.Register(s)
	log.Printf("grpc server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
