from flask import Flask
from flask_restplus import Api, Resource
from webargs.flaskparser import use_args
from webargs import fields

from pymodm.errors import DoesNotExist
from .models import DatasetMeta, DatasetEntry, to_ordered_list

__all__ = ["app"]

app = Flask(__name__)
api = Api(app)


@api.route("/dataset")
class Datasets(Resource):
    def get(self):
        dataset_meta = DatasetMeta.objects.raw({})
        resp = [meta.as_dict() for meta in dataset_meta]
        return {"datasets": resp}, 200

    @use_args({"name": fields.Str(required=True), "keys": fields.List(cls_or_instance=fields.Str, required=True)})
    def post(self, args):
        DatasetMeta(args["name"], keys=args["keys"]).save()
        dataset_meta = DatasetMeta.objects.raw({"_id": args["name"]}).first()
        return {args["name"]: dataset_meta.as_dict()}, 200


@api.route("/dataset/<string:name>")
class Dataset(Resource):
    def get(self, name):
        try:
            dataset_meta = DatasetMeta.objects.raw({"_id": name})
            target_meta = dataset_meta.first().as_dict()
            return target_meta, 200
        except DoesNotExist:
            return {"error": "no such dataset"}, 404

    @use_args({"entry": fields.Dict(keys=fields.Str, values=fields.Float, required=True)})
    def post(self, args, name):
        try:
            dataset_meta = DatasetMeta.objects.raw({"_id": name})
            target_meta = dataset_meta.first().as_dict()
            if set(target_meta["keys"]) != set(args["entry"].keys()):
                # The keys passed as input via the dataset API are different from the ones
                # used in the dataset.
                return {"error": "wrong keys for the specified dataset"}, 422
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
