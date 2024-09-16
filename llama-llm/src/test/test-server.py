import grpc
from concurrent import futures
import time

import sys
import os

sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '../operations_api/')))

from rpc import marble_pb2
from rpc import marble_pb2_grpc

class OperationTrackerServicer(marble_pb2_grpc.OperationTrackerServicer):
    def StartOperation(self, request, context):
        print(f"Received request to start {request.nOperations} operations")
        return marble_pb2.StartOperationResp(ok=True)

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    marble_pb2_grpc.add_OperationTrackerServicer_to_server(OperationTrackerServicer(), server)
    server.add_insecure_port('[::]:50052')
    server.start()
    print("Server started on port 50052")
    try:
        while True:
            time.sleep(86400)
    except KeyboardInterrupt:
        server.stop(0)

if __name__ == '__main__':
    serve()
