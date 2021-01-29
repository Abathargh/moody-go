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

check_running() {
  for pfile in bin/tmp/*.pid; do
    [ -f "$pfile" ] && echo "An instance of the app is already running, close with ./moody-hub.sh -c"; exit 1;
  done
}


close() {
  for pid_file in bin/tmp/*; do
    kill -INT "$(cat "$pid_file")" && rm "$pid_file"
  done
}


while getopts ":hsc" opt; do
  case $opt in
    h ) usage; exit 0;;
    s ) background=1;;
    c ) close; exit 0;;
    : ) echo "Missing argument for option -$OPTARG"; usage; exit 1;;
    \?) echo "Unknown option -$OPTARG"; usage; exit 1;;
  esac
done

# check for build/run dependencies

check_running

mosq_check=$(command -v mosquitto)
[ -z "$mosq_check" ] && echo "Missing dependency: mosquitto"; exit 1;

go_check=$(command -v go)
[ -z "$go_check" ] && echo "Missing dependency: golang-go"; exit 1;

npm_check=$(command -v npm)
[ -z "$npm_check" ] && echo "Missing dependency: npm"; exit 1;


# build the bin/ folder is not present

[ -d bin/ ] || echo "Building the moody front architecture"; make build-front
[ -z "$(command -v serve)" ] && npm install -g serve

# run and save pid to file for a later close

mkdir -p bin/log
mkdir -p bin/tmp

mosquitto -c ./broker/mosquitto.conf -v > bin/log/mosquitto.log 2> bin/log/mosquitto.log &
echo $! > bin/tmp/mosquitto.pid

serve -l 3000 -s bin/build &
echo $! > bin/tmp/webapp.pid

if [ -z $background ]; then
  GODEBUG="x509ignoreCN=0" ./bin/gateway 2>&1 | tee -a bin/log/gateway.log
  close
else
  GODEBUG="x509ignoreCN=0" ./bin/gateway 2>&1 | tee -a bin/log/gateway.log &
  echo $! > bin/tmp/gateway.pid
fi
