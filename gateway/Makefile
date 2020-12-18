gateway :
	go mod download
	go build -o gateway .

start : gateway
	mkdir -p log
	./gateway > log/out.log 2>>log/err.log & \
	echo $$! > pid.tmp

stop :
	-(kill -INT $$(cat pid.tmp)) && rm pid.tmp

clean :
	-rm gateway

.PHONY : start
.PHONY : stop
.PHONY : clean
