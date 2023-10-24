/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package grpc_server

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"grpc-zmq-sse/db"
	pb "grpc-zmq-sse/generated-proto"
	zmq_local "grpc-zmq-sse/zmq-local"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("GRPC Server received: %v", in.GetName())
	zmq_local.GlobalPublisher.Send(in.GetName(), 0)
	log.Println("ZMQ PUB Sent: " + in.GetName())

	err := db.GlobalConnection.Create(&db.Dump{Message: in.GetName()}).Error
	if err != nil {
		fmt.Printf("Error: %s", err)
	}

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
