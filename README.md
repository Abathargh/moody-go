## **moody-go**
This is a tentative port of the Moody project (https://github.com/Antimait/Moody) to an infrastucture written in go.
The main reason behind this port is to fix a big issue stemming from two conflicting libraries in use in the python3 
version (eventlet/flask-socketio with threading/multiprocessing).

Install via:
```bash
go get github.com/Abathargh/moody-go
cd $GOPATH/src/github.com/Abathargh/moody-go
./build.sh

# Or use our docker image (supporting x-64, arm32v7 and arm64v8 arch)
docker pull abathargh/moody-go
docker run --name moody-go --net=host -p 1883:1883 abathargh/moody-go
```

You may have to install the sqlite3 drivers/g++ to build on your local machine