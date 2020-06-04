import pymodm


class DatasetMeta(pymodm.MongoModel):
    name = pymodm.fields.CharField(primary_key=True)
    keys = pymodm.fields.ListField()

    def as_dict(self):
        return {"name": self.name, "keys": self.keys}


class DatasetEntry(pymodm.MongoModel):
    dataset = pymodm.fields.ReferenceField(DatasetMeta, on_delete=pymodm.ReferenceField.CASCADE)
    entry = pymodm.fields.ListField()
