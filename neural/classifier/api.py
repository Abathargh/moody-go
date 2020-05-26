from xmlrpc.server import SimpleXMLRPCServer, SimpleXMLRPCRequestHandler
from neural.classifier.neural import NeuralInterface

RPC_PORT = 9998


class RequestHandler(SimpleXMLRPCRequestHandler):
    rpc_paths = ("/neural",)


def serve_rpc():
    with SimpleXMLRPCServer(("localhost", RPC_PORT), requestHandler=RequestHandler) as server:
        server.register_introspection_functions()
        neural_network = NeuralInterface()

        server.register_function(neural_network.train, "train")
        server.register_function(neural_network.predict, "predict")

        server.serve_forever()
