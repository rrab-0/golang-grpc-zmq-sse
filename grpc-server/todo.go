package grpc_server

import (
	"context"
	"flag"
	"fmt"
	"grpc-zmq-sse/db"
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
	// completed := strconv.FormatBool(todo.GetCompleted())

	// fmt.Println("=======================================")
	// fmt.Println(
	// 	" ID: "+todo.GetId()+"\n",
	// 	"Title: "+todo.GetTitle()+"\n",
	// 	"Description: "+todo.GetDescription()+"\n",
	// 	"isCompleted: "+completed+"\n",
	// 	"CreatedAt: "+todo.GetCreatedAt().String()+"\n",
	// 	"CreatedAt_formatted: "+todo.GetCreatedAt().AsTime().Format("2006-01-02 15:04:05")+"\n",
	// 	"UpdatedAt: "+todo.GetUpdatedAt().String()+"\n",
	// 	"UpdatedAt_formatted: "+todo.GetCreatedAt().AsTime().Format("2006-01-02 15:04:05")+"\n",
	// )
	// fmt.Println("=======================================")

	dbTodo := &db.Todo{}
	dbTodo.Title = todo.GetTitle()
	dbTodo.Description = todo.GetDescription()
	if err := db.GlobalConnection.Create(&dbTodo).Error; err != nil {
		return nil, err
	}

	return &pbTodo.CreateTodoResponse{Id: todo.GetId()}, nil
}

func (ts *todoServer) GetTodo(ctx context.Context, in *pbTodo.GetTodoRequest) (*pbTodo.GetTodoResponse, error) {
	todoId := in.GetId()

	fmt.Println("=======================================")
	fmt.Println(" ID: " + todoId)
	fmt.Println("=======================================")

	var dbTodo db.Todo
	if err := db.GlobalConnection.First(&dbTodo, "id = ?", todoId).Error; err != nil {
		return nil, err
	}

	return &pbTodo.GetTodoResponse{Activity: &pbTodo.Todo{
		Id:          dbTodo.ID.String(),
		Title:       dbTodo.Title,
		Description: dbTodo.Description,
		Completed:   dbTodo.Completed,
	}}, nil
}

// func (ts *todoServer) ListTodo(ctx context.Context, in *pbTodo.ListTodoRequest) (*pbTodo.ListTodoResponse, error) {
// 	limit := in.GetLimit()
// 	not_completed := in.GetNotCompleted()

// 	fmt.Println("=======================================")
// 	fmt.Println(" List Todo")
// 	fmt.Println("=======================================")

// 	var todos []db.Todo
// 	if err := db.GlobalConnection.Find(&todos).Error; err != nil {
// 		return nil, err
// 	}

// 	return &pbTodo.ListTodoResponse{
// 		Activities: []*pbTodo.Todo{
// 			 &pbTodo.Todo{},
// 			 &pbTodo.Todo{},
// 		},
// 	}, nil
// }

func (ts *todoServer) DeleteTodo(ctx context.Context, in *pbTodo.DeleteTodoRequest) (*pbTodo.DeleteTodoResponse, error) {
	todoId := in.GetId()

	fmt.Println("=======================================")
	fmt.Println(" Delete Todo")
	fmt.Println(" ID: " + todoId)
	fmt.Println("=======================================")

	var todo db.Todo
	if err := db.GlobalConnection.Where("id = ?", todoId).Delete(&todo).Error; err != nil {
		return nil, err
	}

	return &pbTodo.DeleteTodoResponse{Id: todoId}, nil
}

func (ts *todoServer) UpdateTodo(ctx context.Context, in *pbTodo.UpdateTodoRequest) (*pbTodo.UpdateTodoResponse, error) {
	todoId := in.GetActivity().GetId()

	fmt.Println("=======================================")
	fmt.Println(" Update Todo")
	fmt.Println(" ID: " + todoId)
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
