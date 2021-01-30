# Builds and runs the moody architecture
# This was tested on debian-like linuxes; the following programs are required:
#     - make
#     - mosquitto
#     - node + npm
#     - golang (>= 1.14)
#
# On a debian-like linux, you can just:
#     sudo apt install golang-go nodejs npm mosquitto
#

usage() {
  printf "Starts the moody hub/gw; mosquitto, node+npm and go are required to build and run the application\n"
  printf "\t[-h] Display this message\n"
  printf "\t[-s] Starts the hub in background\n"
}

bin_config() {
  custom_mosq='autosave_interval 1800
persistence true
retained_persistence true
persistence_file m2.db
persistence_location '$1'
connection_messages true
log_timestamp true


listener 8883
cafile '$1'/ca.crt
certfile '$1'/server.crt
keyfile '$1'/server.key'
  echo "$custom_mosq" > "$1/mosquitto.conf"
}

close_all() {
  kill "$(cat bin/tmp/to-close)"
}

check_running() {
  shopt -s nullglob
  for pfile in bin/tmp/*; do
    [ -f "$pfile" ] && echo "An instance of the app is already running, close with ./moody-hub.sh -c"; exit 1;
  done
}


while getopts ":hsc" opt; do
  case $opt in
    h ) usage; exit 0;;
    s ) background=1;;
    c ) close_all; exit 0;;
    \?) echo "Unknown option -$OPTARG"; usage; exit 1;;
  esac
done

# check for build/run dependencies
mkdir -p bin/log
mkdir -p bin/tmp

check_running

mosq_check=$(command -v mosquitto)
[ -z "$mosq_check" ] && echo "Missing dependency: mosquitto" && exit 1;

go_check=$(command -v go)
[ -z "$go_check" ] && echo "Missing dependency: golang-go" && exit 1;

npm_check=$(command -v npm)
[ -z "$npm_check" ] && echo "Missing dependency: npm" && exit 1;


# build the bin/ folder is not present

[ -d bin/ ] || (echo "Building the moody front architecture" && make build-front)
[ -z "$(command -v serve)" ] && npm install -g serve
[ -d bin/broker ] || (mkdir -p bin/broker && cp broker/* bin/broker && bin_config "$PWD/bin/broker")

# run everything and save logs

trap 'kill $(jobs -p)' SIGINT
mosquitto -c bin/broker/mosquitto.conf -v 2>&1 | tee -a log/mosquitto.log &
serve -l 3000 -s bin/build > /dev/null 2> /dev/null &

if [ -z $background ]; then
  (cd bin || exit 1; GODEBUG="x509ignoreCN=0" ./gateway 2>&1 | tee -a log/gateway.log)
  echo "Bye!"
else
  (cd bin || exit 1; GODEBUG="x509ignoreCN=0" ./gateway 2>&1 | tee -a log/gateway.log &)
  jobs -p > bin/tmp/to-close
fi
