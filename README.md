## **moody-go**
This is a tentative port of the Moody project (https://github.com/Antimait/Moody) to an infrastucture written in go.
The main reason behind this project is to fix a big issue stemming from two conflicting libraries in use in the python3 
version (eventlet/flask-socketio with threading/multiprocessing).

To install clone the repo inside your GOPATH, build it and run:
```bash
git clone https://github.com/abathargh/moody-go
go install

# Or build a docker image
docker build --rm . -t moody-server
docker run --rm --name moody-server --net=host -p 1883:1883 moody-server
```
