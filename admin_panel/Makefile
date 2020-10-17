include .env

node_modules : package.json package-lock.json
	npm install

dev : node_modules
	npm start

build : node_modules
	npm run build
	npm install serve

start :
	mkdir -p log
	node_modules/.bin/serve -l tcp://0.0.0.0:${WEBAPP_PORT} -s build/ > log/out.log 2>>log/err.log & \
	echo $$! > pid.tmp

stop :
	-(kill -INT $$(cat pid.tmp)) && rm pid.tmp

clean :
	rm -rf node_modules
	rm -rf build

.PHONY : dev
.PHONY : start
.PHONY : stop
.PHONY : clean
