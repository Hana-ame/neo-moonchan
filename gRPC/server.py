from concurrent import futures
import grpc
import greeter_pb2
import greeter_pb2_grpc
import calculator_pb2
import calculator_pb2_grpc


class GreeterServicer(greeter_pb2_grpc.GreeterServicer):
    def SayHello(self, request, context):
        return greeter_pb2.HelloReply(
            message=f'Hello, {request.name}!'
        )
    
class CalculatorServicer(calculator_pb2_grpc.CalculatorService):
    def Add(self, request, context):
        print(request, context)
        return calculator_pb2.AddResponse(
            result = request.a + request.b
        )

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=128 ))
    greeter_pb2_grpc.add_GreeterServicer_to_server(GreeterServicer(), server)
    calculator_pb2_grpc.add_CalculatorServiceServicer_to_server(CalculatorServicer(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    print("Server running on port 50051...")
    server.wait_for_termination()

if __name__ == '__main__':
    serve()