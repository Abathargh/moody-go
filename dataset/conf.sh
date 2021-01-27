#!/usr/bin/env bash

while getopts ":h:p:n:" opt; do
  case $opt in
    h ) db_host=$OPTARG;;
    p ) db_port=$OPTARG;;
    n ) db_name=$OPTARG;;
    : ) echo "Missing argument for option -$OPTARG"; exit 1;;
    \?) echo "Unknown option -$OPTARG"; exit 1;;
  esac
done

[ -z "$db_host" ] && db_host="moody-dataset-db"
[ -z "$db_port" ] && db_port="27017"
[ -z "$db_name" ] && db_name="dataset"


conf='{
  "db_host": "'$db_host'",
  "db_port": '$db_port',
  "db_name": "'$db_name'"
}'

mkdir -p data
echo "$conf" > data/conf.json