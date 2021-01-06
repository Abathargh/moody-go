# *gateway*

The gateway receives data from sensors and forwards it to the services and to actuators, depending on the settings and its internal state.

## Contents
- [Requirements](#requirements)
- [Installation](#installation)
    - [As a standalone application](#as-a-standalone-application)
    - [As part of the moody front architecture](#as-a-part-of-the-moody-front-architecture)
- [Functionalities](#functionalities)
    - [Services](#services)
    - [Actuator Mode](#actuator-mode)
    - [Neural API](#neural-api)

## Requirements

If you want to run the gateway directly on your machine, you will need to install go (>= 1.13) and make (optional); the app was developed with reference to the golang-go version.

You can also run it through docker if you don't want to install go.

## Installation

Before proceeding with the following instructions, remember to generate the configuration file and the certificates, 
following the steps contained in the root directory readme.

### As a standalone application
You can run the gateway as a standalone application, but you will need to pass the address of the broker you're using, and the one exposing the services API. 

```bash
# directly build the binaries, you will need to install go on your machine

# using make
make build
make start

# close with
make stop

# doing everything manually
cd gateway
go mod download
go build -o moody-gateway .

# or use the prebuilt docker image
docker run --name moody-gateway -v ./gateway/data:/data -p 7000:7000 abathargh/moody-go-gateway:latest
```

### As a part of the moody front architecture

You can use the moody-backend.yml compose configuration file to set up the backend on a remote machine and pass its address to the gateway. 

Another way to deploy the gateway is to just run compose using moody-gw.yml, that will set up the whole front side of the moody application, consisting of the broker, the gateway and the admin panel.

```bash
# From the root directory of the project
# On a remote machine
docker-compose -f moody-backend.yml --build -d

# On a local machine
docker-compose -f moody-gw.yml up --build -d
```

The pre-built images shared on docker hub and used in the compose files and in the examples are compatible with the following architectures:

- amd64
- arm32v7
- aarch64 (arm64 or arm64v8)

## Functionalities

The application functions as a gateway for the requests that have to be forwarded to the upper layer of the architecture (dataset, neural services, etc.) and as a aggregator of data coming from the sensors. 

It's used as both a middle node for the underlying WSAN and as an entry node to forward requests to the services by the admin panel app.


### Services
A Moody sensor forwards data through an MQTT topic that has a one to one correspondence with a service. When a service is activated, messages incoming on its topic are accepted and saved onto a table in memory: when new occurrenced of that kind of data are received, they are forwarded to the next layer in the architecture.

Data incoming from activated is exposed through the socketio API so that apps can capture it in real time.

### Actuator Mode

Actuators can be in an actuation mode, where they work with reference to the situations received, or they can be in a server mode in which their internal mappings can be modified.

The global state of the actuators is managed by the gateway; when switching to server mode, every actuator node is notified and keeps being notified to be sure that new nodes switch to the correct mode. By default the actuators are in actuating mode.

### Neural API

The collected data coming from sensors is routed by taking in consideration the neural state of the app, which can be stopped, collecting or predicting.

When creating, and then using, a dataset, the active services decide which keys are to be used in the dataset.

When collecting, it's important to set the situation that is tied to the data currently being collected. This is key not only to correctly populate the datasets, but also to make the collecting process begin, since the application will ignore data not bound by a situation.

When predicting, a snapshot of the data obtained through the active services is obtained and forwarded to the neural service; the result is then broadcasted to the actuators.
