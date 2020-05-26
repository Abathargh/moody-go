## **moody-go**
### *gateway*

This directory contains a dockerized version of the gateway services of the application. You can build it for your 
target CPU with one of the available Dockerfiles or:
- Build it via go:
```bash
go get github.com/Abathargh/moody-go/gateway
cd $GOPATH/src/github.com/Abathargh/moody-go/gateway
./build.sh
```

- Pull the image via:
```bash
docker pull abathargh/moody-go
```