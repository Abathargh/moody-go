# *dataset*

The implementation of the dataset service of the moody project; a very simple service which stores and analyzes the 
collected data.

## Contents
- [Requirements](#requirements)
- [Installation](#installation)
    - [As a standalone application](#as-a-standalone-application)
    - [As part of the moody front architecture](#as-part-of-the-moody-back-architecture)
    - [Pre-built wheels for ARM architectures](#pre-built-wheels-for-arm-architectures)

## Requirements

This module has been written, tested and deployed using python 3.8, it should work with any version of python >= 3.5.
All the dependencies are in the requirements.txt file.

## Installation

### As a standalone application

You can run this module as a flask app for non-production purposes.
Otherwise, run it through gunicorn + nginx. In either case, you will need an instance of mongodb running, since this module interfaces with it.

```bash
# create and activate a virtual env, then...
pip3 install -r requirements.txt

# run via flask
FLASK_APP=app.py flask run

# run via gunicorn
gunicorn -b 0.0.0.0:80 app:app
```

You can also run the module using our docker image:

```bash
docker run --name moody-dataset -p <port>:80 abathargh/moody-go-dataset:latest
```

### As part of the moody back architecture

```bash
# From the root directory of the project
# On a remote machine
docker-compose -f moody-backend.yml --build -d
```

### Pre-built wheels for ARM architectures

If you use our pre-built image, or you're going to build an image from the Dockerfile in this module, know that you're using some pre-built python3 ARM wheels from my python3-arm-wheels repo. 
These were built from source and were not modified by any means, but, as specified in the repo license, I have no responsibilities for anything that may happen by using them. 

If you don't trust those wheels, it's better to install the libraries from source yourself, or to build your own wheels.
