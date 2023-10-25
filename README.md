# Notes:

- ignore greeter_client/server, it is not used
- data flows from top to bottom:
    - grpc client (postman/browser)
    - grpc server
    - zmq pub
    - zmq sub
    - sse server
    - sse client (postman/browser)