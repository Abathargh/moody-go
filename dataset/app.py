# This app is started through gunicorn, hence the missing main
import logging
import socket
import pymodm
import json
import sys

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


# Wait for the db service to be up and running

attempt = 0
timeout = 15
max_attempts = 5
checking_db = True

while checking_db:
    try:
        with socket.create_connection((db_host, db_port), timeout=timeout):
            checking_db = False
    except OSError:
        attempt += 1
        if attempt > max_attempts:
            logging.error("The database service is unreachable")
            sys.exit(1)


pymodm.connect("mongodb://{}:{}/{}".format(db_host, db_port, db_name))
logging.info("Connected to mongodb")
