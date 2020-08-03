from typing import Set, List
import pymodm


neural_meta_keys = {"situation", "hour", "minute"}


def strip_neural_meta(keys: List[str]) -> Set[str]:
    """
    Strips the neural meta key from the entry, returning a copy of the dict
    without the keys
    """
    return set(keys) - neural_meta_keys


class DatasetMeta(pymodm.MongoModel):
    name = pymodm.fields.CharField(primary_key=True)
    keys = pymodm.fields.ListField()

    def as_dict(self):
        return {"name": self.name, "keys": self.keys}


class DatasetEntry(pymodm.MongoModel):
    dataset = pymodm.fields.ReferenceField(
        DatasetMeta, on_delete=pymodm.ReferenceField.CASCADE)
    entry = pymodm.fields.ListField()
