from webargs.flaskparser import use_args, parser, abort
from webargs import fields

from flask_restplus import Api, Resource
from flask import Flask

from pymodm.errors import DoesNotExist
import datetime
import logging

from typing import List
from .models import DatasetMeta, DatasetEntry, strip_neural_meta, to_ordered_list
from .neural import NeuralInterface

__all__ = ["app"]

app = Flask(__name__)
api = Api(app)

classifier = NeuralInterface()
# Maybe a cache with the last n predictions?


def get_dataset(dataset_name: str) -> List[List[float]]:
    logging.info("Retrieving dataset {} from the db".format(dataset_name))
    entries = DatasetEntry.objects.raw({"dataset": dataset_name}).only("entry")
    return [entry.entry for entry in entries]


@api.route("/predict")
class Prediction(Resource):
    @use_args({
        "dataset": fields.Str(required=True),
        "query": fields.Dict(keys=fields.Str, values=fields.Float, required=True)
    })
    def post(self, args):
        try:
            dataset_meta = DatasetMeta.objects.raw({"_id": args["dataset"]})
            target_meta = dataset_meta.first().as_dict()
            if strip_neural_meta(target_meta["keys"]) != set(args["query"].keys()):
                return {"error": "wrong keys for the specified dataset"}, 422
            if args["dataset"] != classifier.dataset:
                data = get_dataset(args["dataset"])
                logging.info("Starting the training session with dataset {}".format(args["dataset"]))
                classifier.train(args["dataset"], list(target_meta["keys"]), data)
            stripped_keys = list(strip_neural_meta(target_meta["keys"]))
            ordered_query = to_ordered_list(stripped_keys, args["query"])
            return {"situation": int(classifier.predict(ordered_query))}, 200
        except DoesNotExist:
            return {"error": "no such dataset"}, 404


@api.route("/data")
class Datasets(Resource):
    def get(self):
        dataset_meta = DatasetMeta.objects.raw({})
        resp = [meta.as_dict() for meta in dataset_meta]
        return {"datasets": resp}, 200

    @use_args({"name": fields.Str(required=True), "keys": fields.List(cls_or_instance=fields.Str, required=True)})
    def post(self, args):
        # Decorate the incoming data with information about the situation and the time
        args["keys"].append("situation")
        args["keys"].append("hour")
        args["keys"].append("minute")
        DatasetMeta(args["name"], keys=args["keys"]).save()
        dataset_meta = DatasetMeta.objects.raw({"_id": args["name"]}).first()
        return dataset_meta.as_dict(), 200


@api.route("/data/<string:name>")
class Dataset(Resource):
    def get(self, name):
        try:
            dataset_meta = DatasetMeta.objects.raw({"_id": name})
            target_meta = dataset_meta.first().as_dict()
            return target_meta, 200
        except DoesNotExist:
            return {"error": "no such dataset"}, 404

    @use_args({"situation": fields.Integer(required=True), "entry": fields.Dict(keys=fields.Str, values=fields.Float, required=True)})
    def post(self, args, name):
        try:
            dataset_meta = DatasetMeta.objects.raw({"_id": name})
            target_meta = dataset_meta.first().as_dict()

            if strip_neural_meta(target_meta["keys"]) != set(args["entry"].keys()):
                # The keys passed as input via the dataset API are different from the ones
                # used in the dataset.
                return {"error": "wrong keys for the specified dataset"}, 422

            # Decorate the incoming data with information about the situation and the time
            now = datetime.datetime.now()
            args["entry"]["situation"] = args["situation"]
            args["entry"]["hour"] = now.hour
            args["entry"]["minute"] = now.minute
            DatasetEntry(dataset=name, entry=to_ordered_list(target_meta["keys"], args["entry"])).save()
            return {"dataset": name, "entry": args["entry"]}, 200
        except DoesNotExist:
            return {"error": "no such dataset"}, 404

    def delete(self, name):
        try:
            dataset_meta = DatasetMeta.objects.raw({"_id": name})
            deleted_meta = dataset_meta.first().as_dict()
            dataset_meta.delete()
            return deleted_meta, 200
        except DoesNotExist:
            return {"error": "no such dataset"}, 404


@parser.error_handler
def handle_request_parsing_error(err, req, schema, *, error_status_code, error_headers):
    logging.error(err, error_headers)
    abort(error_status_code, errors=err.messages)
