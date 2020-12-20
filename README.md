# *moody-go*
This is a tentative port of the Moody project (https://github.com/Antimait/Moody) to a microservice-oriented infrastructure written in go.
The main reason behind this port is to fix a big issue stemming from two conflicting libraries in use in the python3 
version (eventlet/flask-socketio with threading/multiprocessing).

The actual implementation contains many services running on a series of docker containers, with the original 
python3/moody neural service.

## Contents

- [Installation](#installation)
    - [Configuration and certificates setup](#configuration-and-certificates-setup)
    - [Run via docker-compose](#run-via-docker-compose)
## Installation

Clone the repo and cd into it:

```bash
git clone https://github.com/Abathargh/moody-go
cd moody-go
```


### Configuration and certificates setup

The following will allow you to set up the configuration files and the certificates for the broker and gateway.

When running **conf_init.sh** you will be prompted with a request to generate a default json configuration file or 
a custom one, which will be guided. The generated file can be found in **gateway/data/conf.json** and can be 
modified manually or by running the script again. You can also directly use the script to generate a default file.

```bash
# This will guide you through the procedure
./conf_init.sh

# This will automatically generate a default configuration file
./conf_init.sh default 
```

The **ca_gen.sh** script will generate eveything that's needed to enable the usage of MQTT with TLS. If it's called
without arguments, it will generate the certificates based on the hostname of the machine that you're running the script 
on. If you want to generate certificates for another machine to use, you can call the script by passing the hostname as 
the only argument.

```bash
# Generates certificates for the current machine, via $HOSTNAME
./ca_gen.sh

# Generates certificates for the machine identified by hostnmae
./ca_gen.sh hostname
```

### Run via docker-compose:

Once you have initialized the configuration files and the certificates, you can start using moody through docker,
interfacing with the admin panel reachable from http://localhost:3000.

```bash
docker-compose up --build -d
```


Pre-built images for each service are available at https://hub.docker.com/u/abathargh.
More instructions about every feature can be found in each subfolder.


### Run front-side via make

If your backend services are hosted on another machine, where the api gateway is reachable, you can build and run 
the front-side (broker + gateway + webapp) via make. In this case you will need to install nodejs, npm mosquitto and 
golang, then run:

```bash
make build-front
make run-front

# Stop with
make stop-front
```