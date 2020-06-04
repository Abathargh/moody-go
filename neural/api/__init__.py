import logging
from flask import Flask
from flask_restplus import Api, Resource
from webargs.flaskparser import use_args
from webargs import fields

from typing import List
from pymodm.errors import DoesNotExist
from .models import DatasetMeta, DatasetEntry
from .neural import NeuralInterface

__all__ = ["app"]

app = Flask(__name__)
api = Api(app)

classifier = NeuralInterface()
# Maybe a cache with the last n predictions?


def get_dataset(dataset_name: str) -> List[List[float]]:
    logging.info("Retrieving dataset {} from the db".format(dataset_name))
    entries = DatasetEntry.objects.raw({"dataset", dataset_name}).only("entry")
    return [entry for entry in entries]


@api.route("/predict")
class Prediction(Resource):
    @use_args({
        "dataset": fields.Str(required=True),
        "query": fields.Dict(keys=fields.Str, values=fields.Float, required=True)
    })
    def post(self, args):
        try:
            dataset_meta = DatasetMeta.objects.raw({"_id": name})
            target_meta = dataset_meta.first().as_dict()
            if set(target_meta["keys"]) != set(args["query"].keys()):
                return {"error": "wrong keys for the specified dataset"}, 422
            if args["dataset"] != classifier.dataset:
                data = get_dataset(args["dataset"])
                logging.info("Starting the training session with dataset {}".format(dataset_name))
                classifier.train(args["dataset"], list(args["query"].keys()), data)
            return {"situation": classifier.predict(args["query"])}, 200
        except DoesNotExist:
            return {"error": "no such dataset"}, 404
