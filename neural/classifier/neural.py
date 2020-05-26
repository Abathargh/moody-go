"""
Neural mode exposing the interface to the neural network.
A global instance of the neural network ensures the single interface and is
accessible in a controlled manner by means of the get_neural function.
"""

import numpy
import logging
import datetime
import pandas as pd

from pathlib import Path
from typing import Union, List
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import StandardScaler
from sklearn.neural_network import MLPClassifier

logger = logging.getLogger(__name__)

MOODY_CONF = ".moody"
MOODY_PATH = Path(Path.home(), MOODY_CONF)


class NeuralInterface:
    """
    A class that offers an interface with the neural network.
    Implemented here as a class to support future use with multiple instances for different uses.
    """

    TIME_LIMIT = 3600  # Seconds, one update per hour
    # Fine tuned parameters, to be changed accordingly to the specific instance of neural network
    # in the query release, then start the neural training
    HIDDEN_LAYER_SIZE = (20, 30, 20)  # 3 layer of hidden nodes, with 20, 30, 20 nodes each
    MAX_ITER = 500  # epochs

    def __init__(self):
        self._dataset = None
        self._neural_net = None
        self._datatypes = None
        self._scaler = None

    def train(self, dataset_name: str, datatypes: List[str]) -> bool:
        """
        Initializes the data collected to start predicting using the data set with session_id id
        :dataset_id: int, the id of thee training session to use for the predictions
        :return: bool, representing the success of the training
        """
        datatypes = datatypes
        datatypes.append("hour")
        datatypes.append("minute")

        try:
            with open(Path(MOODY_PATH, dataset_name), "r") as file:
                lines = [line for line in file.read().split("\n")]
                session_data = [[char for char in line.split(",")] for line in lines]
        except FileNotFoundError:
            # Dataset doesn't exists
            raise
        else:
            neural_data = pd.DataFrame(numpy.array(session_data), columns=numpy.array(datatypes))
            datatypes.remove("situation")

            data = neural_data.drop("situation", axis=1)
            situations = neural_data["situation"]
            train_data, test_data, train_situations, test_situations = train_test_split(data, situations)

            # Format, test and transform the data accordingly to the neural network
            self._scaler = StandardScaler()
            self._scaler.fit(train_data)
            train_data = self._scaler.transform(train_data)
            test_data = self._scaler.transform(test_data)

            if not self._neural_net:
                self._neural_net = MLPClassifier(hidden_layer_sizes=self.HIDDEN_LAYER_SIZE, max_iter=self.MAX_ITER)
            self._neural_net.fit(train_data, train_situations)

            logger.info("training accuracy: {:3f}".
                        format(self._neural_net.score(train_data, train_situations)))
            logger.info("test accuracy: {:3f}".
                        format(self._neural_net.score(test_data, test_situations)))

            self._datatypes = datatypes
            self._dataset = dataset_name

            return True

    def predict(self, query: List[Union[int, float]]) -> int:
        """
        Infer the situation starting from the data just read and the dataset
        obtained through the training phase.

        :param query: the data that has been read that requests the prediction
        :return: int, the id of the inferred situation
        """
        now = datetime.datetime.now()
        hours = int(now.strftime("%H"))
        minutes = int(now.strftime("%M"))
        query.append(hours)
        query.append(minutes)

        query_df = pd.DataFrame([query], columns=self._datatypes)
        query_df = self._scaler.transform(query_df)
        result = self._neural_net.predict(query_df)
        return int(result[0])
