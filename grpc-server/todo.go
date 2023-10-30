package grpc_server

import (
	"context"
	"grpc-zmq-sse/db"
	pbTodo "grpc-zmq-sse/generated-proto-todo"
	zmq_local "grpc-zmq-sse/zmq-local"
	"log"
	"net"
	"os"
	"time"

	zmq "github.com/pebbe/zmq4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type todoServer struct {
	pbTodo.UnimplementedTodoServiceServer
}

func (ts *todoServer) CreateTodo(in *pbTodo.CreateTodoRequest, stream pbTodo.TodoService_CreateTodoServer) error {
	requestTime := time.Now()
	todo := in.GetActivity()
	var dbTodo db.Todo

	dbTodo.Title = todo.GetTitle()
	dbTodo.Description = todo.GetDescription()
	if err := db.GlobalConnection.Create(&dbTodo).Error; err != nil {
		return err
	}

	jsonMsg := "{\"status\":\"created\",\"id\":\"" + dbTodo.ID.String() + "\"}"
	_, err := zmq_local.GlobalPublisher.Send(zmq_local.DefaultTopic+" "+jsonMsg, zmq.DONTWAIT)
	if err != nil {
		log.Printf("ZMQ PUB Error: %s\n", err)
		return err
	}

	resp := &pbTodo.CreateTodoResponse{Id: dbTodo.ID.String()}

	sendTime := time.Now()
	if err := stream.Send(resp); err != nil {
		return err
	}

	sendElapsed := time.Since(sendTime)
	log.Printf("Send response took %s", sendElapsed)

	totalElapsed := time.Since(requestTime)
	log.Printf("Total handler time %s", totalElapsed)
	return nil
	// return &pbTodo.CreateTodoResponse{Id: dbTodo.ID.String()}, nil
}

func (ts *todoServer) GetTodo(ctx context.Context, in *pbTodo.GetTodoRequest) (*pbTodo.GetTodoResponse, error) {
	todoId := in.GetId()
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

func (ts *todoServer) ListTodo(ctx context.Context, in *pbTodo.ListTodoRequest) (*pbTodo.ListTodoResponse, error) {
	limit := in.GetLimit()
	not_completed := in.GetNotCompleted()

	var (
		dbTodos   []db.Todo
		todosGrpc []*pbTodo.Todo
	)

	if err := db.GlobalConnection.Limit(int(limit)).Find(&dbTodos, "completed = ?", not_completed).Error; err != nil {
		return nil, err
	}

	for _, todo := range dbTodos {
		newTodo := &pbTodo.Todo{
			Id:          todo.ID.String(),
			Title:       todo.Title,
			Description: todo.Description,
			Completed:   todo.Completed,
		}
		todosGrpc = append(todosGrpc, newTodo)
	}

	return &pbTodo.ListTodoResponse{
		Activities: todosGrpc,
	}, nil
}

func (ts *todoServer) DeleteTodo(ctx context.Context, in *pbTodo.DeleteTodoRequest) (*pbTodo.DeleteTodoResponse, error) {
	todoId := in.GetId()
	var dbTodo db.Todo

	if err := db.GlobalConnection.Where("id = ?", todoId).Delete(&dbTodo).Error; err != nil {
		return nil, err
	}

	jsonMsg := "{\"status\":\"deleted\",\"id\":\"" + todoId + "\"}"
	_, err := zmq_local.GlobalPublisher.Send(zmq_local.DefaultTopic+" "+jsonMsg, zmq.DONTWAIT)
	if err != nil {
		log.Printf("ZMQ PUB Error: %s\n", err)
		return nil, err
	}

	return &pbTodo.DeleteTodoResponse{Id: todoId}, nil
}

func (ts *todoServer) UpdateTodo(ctx context.Context, in *pbTodo.UpdateTodoRequest) (*pbTodo.UpdateTodoResponse, error) {
	todo := in.GetActivity()
	var dbTodo db.Todo

	if err := db.GlobalConnection.First(&dbTodo, "id = ?", todo.GetId()).Error; err != nil {
		return nil, err
	}

	if todo.GetTitle() != "" {
		dbTodo.Title = todo.GetTitle()
	}

	if todo.GetDescription() != "" {
		dbTodo.Description = todo.GetDescription()
	}

	completedString := ""
	if todo.GetCompleted() {
		completedString = "true"
	}
	completedString = "false"

	if completedString != "" {
		dbTodo.Completed = todo.GetCompleted()
	}

	if err := db.GlobalConnection.Where("id = ?", todo.GetId()).Save(&dbTodo).Error; err != nil {
		return nil, err
	}

	jsonMsg := "{\"status\":\"updated\",\"id\":\"" + dbTodo.ID.String() + "\"}"
	_, err := zmq_local.GlobalPublisher.Send(zmq_local.DefaultTopic+" "+jsonMsg, zmq.DONTWAIT)
	if err != nil {
		log.Printf("ZMQ PUB Error: %s\n", err)
		return nil, err
	}

	return &pbTodo.UpdateTodoResponse{Id: dbTodo.ID.String()}, nil
}

func StartTodo() {
	// var (
	// 	portEnvTodo, _ = strconv.Atoi(os.Getenv("GRPC_TODO_PORT"))
	// 	portTodo       = flag.Int("portTodo", portEnvTodo, "The server port")
	// )

	// lis, err := net.Listen("tcp", fmt.Sprintf(os.Getenv("GRPC_TODO_HOST")+":%d", *portTodo))
	lis, err := net.Listen("tcp", os.Getenv("GRPC_TODO_HOST")+":"+os.Getenv("GRPC_TODO_PORT"))
	if err != nil {
		log.Fatalf("grpc-todo failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pbTodo.RegisterTodoServiceServer(s, &todoServer{})
	reflection.Register(s)
	log.Printf("grpc-todo server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("grpc-todo failed to serve: %v", err)
	}
}
