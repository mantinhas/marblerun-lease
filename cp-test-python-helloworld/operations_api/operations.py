import sys
import os

sys.path.append(os.path.abspath((os.path.dirname(__file__))))

import grpc
from rpc import marble_pb2
from rpc import marble_pb2_grpc

def __get_environ_tls():
    MarbleEnvironmentRootCA = "MARBLE_PREDEFINED_ROOT_CA"
    root = str.encode(os.getenv(MarbleEnvironmentRootCA))

    MarbleEnvironmentPrivateKey = "MARBLE_PREDEFINED_PRIVATE_KEY"
    key = str.encode(os.getenv(MarbleEnvironmentPrivateKey))

    MarbleEnvironmentCertificateChain = "MARBLE_PREDEFINED_MARBLE_CERTIFICATE_CHAIN"
    cert = str.encode(os.getenv(MarbleEnvironmentCertificateChain))

    return root, key, cert

def request_operation(amount=1):
    root, key, cert = __get_environ_tls()
    creds = grpc.ssl_channel_credentials(root, key, cert)
    addr = "localhost:50052"
    with grpc.secure_channel(addr, creds) as channel:
        stub = marble_pb2_grpc.OperationTrackerStub(channel)
        response = stub.StartOperation(marble_pb2.StartOperationReq(nOperations=int(amount)))
    return response.ok
