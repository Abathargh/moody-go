#!/usr/bin/env bash

default='{
    "apiGateway": "http://moody-api-gw/",
    "serverPort": ":7000",
    "mqtt": {
        "host": "moody-broker",
        "port": 1883,
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

usage() {
    printf "builds a correct configuration file for the gateway app\n"
    printf "Usage: ./build.sh [default]\n"
    printf "\tdefault: generates a default configuration file without any user interaction\n"
}

init_default() {
    echo "$default" > conf.json
}

build_conf() {
    echo "Do you want a default configuration file to be crated? [y/n]"
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

    custom='{
    "apiGateway": "'$apiGw'",
    "serverPort": ":'$serverPort'",
    "mqtt": {
        "host": "'$brokerAddr'",
        "port": 1883,
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
    echo "$custom" > conf.json
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






