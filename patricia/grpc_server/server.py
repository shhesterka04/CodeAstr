import grpc
from concurrent import futures
import time
import sys
import os
from grpc_reflection.v1alpha import reflection

sys.path.append(os.path.join(os.path.dirname(__file__), '..', 'gen'))

from gen import llm_service_pb2 as pb2
from gen import llm_service_pb2_grpc as pb2_grpc

class LLMServiceServicer(pb2_grpc.LLMServiceServicer):
    def GetLLM(self, request, context):
        response = pb2.GetLLMResponse()
        response.answer = f"Received prompt: {request.prompt}"
        response.raw_json = '{"status": "success"}'
        return response

class HealthCheckServiceServicer(pb2_grpc.HealthCheckServiceServicer):
    def Check(self, request, context):
        response = pb2.HealthCheckResponse()
        response.status = "SERVING"
        return response

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    pb2_grpc.add_LLMServiceServicer_to_server(LLMServiceServicer(), server)
    pb2_grpc.add_HealthCheckServiceServicer_to_server(HealthCheckServiceServicer(), server)

    SERVICE_NAMES = (
        pb2.DESCRIPTOR.services_by_name['LLMService'].full_name,
        pb2.DESCRIPTOR.services_by_name['HealthCheckService'].full_name,
        reflection.SERVICE_NAME,
    )
    reflection.enable_server_reflection(SERVICE_NAMES, server)

    server.add_insecure_port('[::]:50051')
    server.start()
    try:
        while True:
            time.sleep(86400)
    except KeyboardInterrupt:
        server.stop(0)
