package grpc_server

import (
	"flag"
	"fmt"
	pbTodo "grpc-zmq-sse/generated-proto-todo"
	"log"
	"net"
	"os"
	"strconv"

	"google.golang.org/grpc"
)

type todoServer struct {
	pbTodo.UnimplementedTodoServiceServer
}

// func (ts *todoServer) CreateTodo(ctx context.Context, in *pbTodo.CreateTodoRequest) (*pbTodo.CreateTodoResponse, error) {

// 	// err := db.GlobalConnection.Create(&pbTodo.Todo)
// }

func StartTodo() {
	var (
		portEnv, _ = strconv.Atoi(os.Getenv("GRPC_PORT"))
		port       = flag.Int("port", portEnv, "The server port")
	)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("grpc-todo failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pbTodo.RegisterTodoServiceServer(s, &todoServer{})

	log.Printf("grpc-todo  server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("grpc-todo failed to serve: %v", err)
	}
}
