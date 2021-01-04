"""
Neural mode exposing the interface to the dataset network.
A global instance of the dataset network ensures the single interface and is
accessible in a controlled manner by means of the get_neural function.
"""

import numpy
import logging
import pandas as pd

from typing import List
from datetime import datetime

from sklearn.model_selection import train_test_split
from sklearn.preprocessing import StandardScaler
from sklearn.neural_network import MLPClassifier

__all__ = ["NeuralInterface"]

logging.basicConfig(level=logging.INFO)


class NeuralInterface:
    """
    A class that offers an interface with the dataset network.
    Implemented here as a class to support future use with multiple instances for different uses.
    """

    TIME_LIMIT = 3600  # Seconds, one update per hour
    # Fine tuned parameters, to be changed accordingly to the specific instance of dataset network
    # in the query release, then start the dataset training
    # 3 layer of hidden nodes, with 20, 30, 20 nodes each
    HIDDEN_LAYER_SIZE = (20, 30, 20)
    MAX_ITER = 500  # epochs

    def __init__(self):
        self._dataset = None
        self._neural_net = None
        self._datatypes = None
        self._scaler = None

    @property
    def dataset(self):
        return self._dataset

    def train(self, dataset_name: str, datatypes: List[str], dataset: List[List[float]]) -> None:
        """
        Initializes the data collected to start predicting using the data set with session_id id
        :dataset_id: int, the id of thee training session to use for the predictions
        :return: None
        """

        neural_data = pd.DataFrame(numpy.array(
            dataset), columns=numpy.array(datatypes))
        datatypes.remove("situation")

        data = neural_data.drop("situation", axis=1)
        situations = neural_data["situation"]
        train_data, test_data, train_situations, test_situations = train_test_split(
            data, situations)

        # Format, test and transform the data accordingly to the dataset network
        self._scaler = StandardScaler()
        self._scaler.fit(train_data)
        train_data = self._scaler.transform(train_data)
        test_data = self._scaler.transform(test_data)

        if not self._neural_net:
            self._neural_net = MLPClassifier(
                hidden_layer_sizes=self.HIDDEN_LAYER_SIZE, max_iter=self.MAX_ITER)
        self._neural_net.fit(train_data, train_situations)

        logging.info("training accuracy: {:3f}".
                     format(self._neural_net.score(train_data, train_situations)))
        logging.info("test accuracy: {:3f}".
                     format(self._neural_net.score(test_data, test_situations)))

        self._datatypes = datatypes
        self._dataset = dataset_name

    def predict(self, query: List[float]) -> int:
        """
        Infer the situation starting from the data just read and the dataset
        obtained through the training phase.

        :param query: the data that has been read that requests the prediction
        :return: int, the id of the inferred situation
        """
        now = datetime.now()
        query.append(now.hour)
        query.append(now.minute)

        query_df = pd.DataFrame([query], columns=self._datatypes)
        query_df = self._scaler.transform(query_df)
        result = self._neural_net.predict(query_df)
        return int(result[0])
