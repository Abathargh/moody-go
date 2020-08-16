# moody-go
This is a tentative port of the Moody project (https://github.com/Antimait/Moody) to an infrastructure written in go.
The main reason behind this port is to fix a big issue stemming from two conflicting libraries in use in the python3 
version (eventlet/flask-socketio with threading/multiprocessing).

The actual implementation contains many services running on a series of docker containers, with the original 
python3/moody neural service.

### Update (August 2020)

The infrastructure is completely functional, every _basic_ feature has been implemented.


Run via docker-compose:

```bash
git clone https://github.com/Abathargh/moody-go
cd moody-go

docker-compose up --build -d
```

and the open the admin panel reachable from http://localhost:3000.

Pre-built images for each service are available at https://hub.docker.com/u/abathargh.