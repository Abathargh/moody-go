# moody-go
This is a tentative port of the Moody project (https://github.com/Antimait/Moody) to an infrastructure written in go.
The main reason behind this port is to fix a big issue stemming from two conflicting libraries in use in the python3 
version (eventlet/flask-socketio with threading/multiprocessing).

The actual implementation contains many services running on a series of docker containers, with the original 
python3/moody neural service.

Run via docker-compose:

```bash
git clone https://github.com/Abathargh/moody-go
cd moody-go

docker-compose up --build
```