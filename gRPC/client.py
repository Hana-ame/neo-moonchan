import grpc
import greeter_pb2
import greeter_pb2_grpc

def run():
    with grpc.insecure_channel(
        '[::1]:50051',
        options=[
            ('grpc.enable_http_proxy', 0),  # 禁用 HTTP 代理
            ('grpc.http_connect_target', 'direct://'),  # 强制直连
    ]) as channel:
        stub = greeter_pb2_grpc.GreeterStub(channel)
        response = stub.SayHello(greeter_pb2.HelloRequest(name='Alice'))
    print("Response received:", response.message)

if __name__ == '__main__':
    run()