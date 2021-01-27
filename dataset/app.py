# This app is started through gunicorn, hence the missing main
import logging
import pymodm
import json

from api import app


# Load configuration and startup the Flask app
logging.basicConfig(level=logging.INFO)

try:
    with open("data/conf.json", "r") as conf_file:
        conf = json.loads(conf_file.read())
        db_host = conf["db_host"]
        db_port = conf["db_port"]
        db_name = conf["db_name"]
except FileNotFoundError:
    logging.error(
        "Neural configuration file not found, expected {}".format(conf_file))
except KeyError:
    logging.error(
        "There's an error in your dataset.conf syntax, expected server_addr and server_port fields")


pymodm.connect("mongodb://{}:{}/{}".format(db_host, db_port, db_name))
logging.info("Connected to mongodb")
