#!/usr/bin/env bash

usage() {
    printf "generates the certificate authority and the MQTT Broker tls certificate (self-signed mode)\n"
    printf "Usage: ./ca_gen.sh [hostname]\n"
    printf "\thostname: the hostname of the machine where the gateway is going to run\n"
}

certs_setup() {
    mkdir -p broker gateway/data
    openssl req -new -x509 -days 365 -extensions v3_ca -keyout broker/ca.key -out broker/ca.crt \
        -subj "/C=IT/ST=Calabria/L=Rende/O=Antima.it/CN=antima.it"
    openssl genrsa -out broker/server.key 2048
    openssl req -new -out broker/server.csr -key broker/server.key \
        -subj "/C=IT/ST=Calabria/L=Rende/O=Antima.it/OU=Moody MQTT Broker/CN=$1"
    openssl x509 -req -in broker/server.csr -CA broker/ca.crt -CAkey broker/ca.key -CAcreateserial \
        -out broker/server.crt -days 365
    openssl rsa -in broker/server.key -out broker/server.key

    cp broker/ca.crt gateway/data/
    printf "Done\n"
    openssl x509 -in  broker/server.crt -sha1 -noout -fingerprint

}

if [ "$#" -gt 1 ]; then
    echo "wrong number of args"
    usage
    exit
fi

case $# in
    0) certs_setup "$HOSTNAME";;
    1) certs_setup "$1" ;;
esac
