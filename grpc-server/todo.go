package grpc_server

import (
	"context"
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

func (ts *todoServer) CreateTodo(ctx context.Context, in *pbTodo.CreateTodoRequest) (*pbTodo.CreateTodoResponse, error) {
	todo := in.GetActivity()
	completed := strconv.FormatBool(todo.GetCompleted())

	fmt.Println("=======================================")
	fmt.Println(
		" ID: "+todo.GetId()+"\n",
		"Title: "+todo.GetTitle()+"\n",
		"Description: "+todo.GetDescription()+"\n",
		"isCompleted: "+completed+"\n",
		"CreatedAt: "+todo.GetCreatedAt().String()+"\n",
		"CreatedAt_formatted: "+todo.GetCreatedAt().AsTime().Format("2006-01-02 15:04:05")+"\n",
		"UpdatedAt: "+todo.GetUpdatedAt().String()+"\n",
		"UpdatedAt_formatted: "+todo.GetCreatedAt().AsTime().Format("2006-01-02 15:04:05"),
	)
	fmt.Println("=======================================")

	return &pbTodo.CreateTodoResponse{Id: todo.GetId()}, nil
}

func (ts *todoServer) GetTodo(ctx context.Context, in *pbTodo.GetTodoRequest) (*pbTodo.GetTodoResponse, error) {
	todoId := in.GetId()

	fmt.Println("=======================================")
	fmt.Println(" ID: " + todoId)
	fmt.Println("=======================================")

	return &pbTodo.GetTodoResponse{Activity: &pbTodo.Todo{
		Id: todoId,
	}}, nil
}

func (ts *todoServer) ListTodo(ctx context.Context, in *pbTodo.ListTodoRequest) (*pbTodo.ListTodoResponse, error) {
	fmt.Println("=======================================")
	fmt.Println(" List Todo")
	fmt.Println("=======================================")

	return &pbTodo.ListTodoResponse{}, nil
}

func (ts *todoServer) DeleteTodo(ctx context.Context, in *pbTodo.DeleteTodoRequest) (*pbTodo.DeleteTodoResponse, error) {
	todoId := in.GetId()

	fmt.Println("=======================================")
	fmt.Println(" Delete Todo")
	fmt.Println("=======================================")

	return &pbTodo.DeleteTodoResponse{Id: todoId}, nil
}

func (ts *todoServer) UpdateTodo(ctx context.Context, in *pbTodo.UpdateTodoRequest) (*pbTodo.UpdateTodoResponse, error) {
	todoId := in.GetActivity().GetId()

	fmt.Println("=======================================")
	fmt.Println(" Update Todo")
	fmt.Println("=======================================")

	return &pbTodo.UpdateTodoResponse{Id: todoId}, nil
}

func StartTodo() {
	var (
		portEnvTodo, _ = strconv.Atoi(os.Getenv("GRPC_TODO_PORT"))
		portTodo       = flag.Int("portTodo", portEnvTodo, "The server port")
	)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *portTodo))
	if err != nil {
		log.Fatalf("grpc-todo failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pbTodo.RegisterTodoServiceServer(s, &todoServer{})

	log.Printf("grpc-todo server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("grpc-todo failed to serve: %v", err)
	}
}
