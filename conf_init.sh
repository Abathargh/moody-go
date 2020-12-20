#!/usr/bin/env bash

default='{
    "apiGateway": "http://moody-api-gw/",
    "serverPort": ":7000",
    "mqtt": {
        "host": "moody-broker",
        "port": 8883,
        "dataTopic": [
            "moody/service/+",
            "moody/actserver"
        ],
        "pubTopics": [
            "moody/actuator/mode",
            "moody/actuator/situation"
        ]
    }
}'

default_mosq='autosave_interval 1800
persistence true
retained_persistence true
persistence_file m2.db
persistence_location /mosquitto/config/
connection_messages true
log_timestamp true


listener 8883
cafile /mosquitto/config/ca.crt
certfile /mosquitto/config/server.crt
keyfile /mosquitto/config/server.key'

usage() {
    printf "builds a correct configuration file for the gateway and for the broker\n"
    printf "Usage: ./build.sh [default]\n"
    printf "\tdefault: generates a default configuration file without any user interaction\n"
}

init_default() {
    mkdir -p broker gateway/data/
    echo "$default" > gateway/data/conf.json
    echo "$default_mosq" > broker/mosquitto.conf
}

build_conf() {
    echo "Do you want a default configuration file to be created? [y/n]"
    while true; do
        printf "> "
        read -r defconf
        case $defconf in
            y) init_default; exit;;
            n) break;;
            *) echo "Choose one between [y/n]"; continue;;
        esac;
    done

    echo "Insert the address where to find the moody api gateway:"
    printf "> "
    read -r apiGw
    echo "Insert the port for the gateway to listen to:"
    printf "> "
    read -r serverPort
    echo "Insert the address of the mqtt broker:"
    printf "> "
    read -r brokerAddr
    echo "Insert the port of the mqtt broker:"
    printf "> "
    read -r brokerPort

    custom='{
    "apiGateway": "'$apiGw'",
    "serverPort": ":'$serverPort'",
    "mqtt": {
        "host": "'$brokerAddr'",
        "port": '$brokerPort',
        "dataTopic": [
            "moody/service/+",
            "moody/actserver"
        ],
        "pubTopics": [
            "moody/actuator/mode",
            "moody/actuator/situation"
        ]
    }
}'

    custom_mosq='autosave_interval 1800
    persistence true
    retained_persistence true
    persistence_file m2.db
    persistence_location /mosquitto/config/
    connection_messages true
    log_timestamp true


    listener '$brokerPort'
    cafile /mosquitto/config/ca.crt
    certfile /mosquitto/config/server.crt
    keyfile /mosquitto/config/server.key'

    mkdir -p broker gateway/data/
    echo "$custom" > gateway/data/conf.json
    echo "$custom_mosq" > broker/mosquitto.conf
}

if [ "$#" -gt 1 ]; then
    echo "wrong number of args"
    usage
    exit
fi

if [ -f conf.json ]; then
    echo "a conf file already exists"
    exit
fi

case $# in
    0) build_conf ;;
    1) if [ ! "$1" = "default" ]; then echo "expected 'default' arg, got $1"; usage; else init_default; fi ;;
esac
