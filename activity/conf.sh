#!/usr/bin/env bash

while getopts ":v:h:p:u:s:n:" opt; do
  case $opt in
    v ) server_port=$OPTARG;;
    h ) db_host=$OPTARG;;
    p ) db_port=$OPTARG;;
    u ) db_user=$OPTARG;;
    s ) db_pass=$OPTARG;;
    n ) db_name=$OPTARG;;
    : ) echo "Missing argument for option -$OPTARG"; exit 1;;
    \?) echo "Unknown option -$OPTARG"; exit 1;;
  esac
done

[ -z "$server_port" ] && server_port=":80"
[ -z "$db_host" ] && db_host="moody-activity-db"
[ -z "$db_port" ] && db_port="5432"
[ -z "$db_user" ] && db_user="postgres"
[ -z "$db_pass" ] && db_pass="password"
[ -z "$db_name" ] && db_name="moody"


conf='{
  "server_port": "'$server_port'",
  "db_host": "'$db_host'",
  "db_port": '$db_port',
  "db_user": "'$db_user'",
  "db_pass": "'$db_pass'",
  "db_name": "'$db_name'"
}'

mkdir -p data
echo "$conf" > data/conf.json