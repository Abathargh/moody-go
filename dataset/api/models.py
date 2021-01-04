from typing import Set, List, Dict
import pymodm


neural_meta_keys = {"situation", "hour", "minute"}


def strip_neural_meta(keys: List[str]) -> Set[str]:
    """
    Strips the dataset meta key from the entry, returning a copy of the dict
    without the keys
    """
    return set(keys) - neural_meta_keys


def to_ordered_list(ordered_keys: List[str], data: Dict[str, float]) -> List[float]:
    """
    Converts a dict containing a single entry of a dataset into a list with the dict
    values in the right order.
    """
    ordered_list = list()
    for key in ordered_keys:
        ordered_list.append(data[key])
    return ordered_list


class DatasetMeta(pymodm.MongoModel):
    name = pymodm.fields.CharField(primary_key=True)
    keys = pymodm.fields.ListField()

    def as_dict(self):
        return {"name": self.name, "keys": self.keys}


class DatasetEntry(pymodm.MongoModel):
    dataset = pymodm.fields.ReferenceField(
        DatasetMeta, on_delete=pymodm.ReferenceField.CASCADE)
    entry = pymodm.fields.ListField()
