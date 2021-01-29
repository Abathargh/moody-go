build-front :
	(cd gateway; make gateway)
	(cd admin_panel; make build)
	mkdir -p bin/data
	cp -r gateway/data bin/
	mv gateway/gateway bin/
	mv admin_panel/build bin/

clean-front :
	rm -rf bin/

run-front :
	@mosquitto -c ./broker/mosquitto.conf -d &> /dev/null
	@cd gateway && make -s start
	@cd admin_panel && make -s start

stop-front :
	-@(kill -INT $$(ps ax | grep mosquitto | grep -v grep | awk '{print $$1}'))
	@cd gateway && make -s stop
	@cd admin_panel && make -s stop

run-dev :
	@mosquitto -c ./broker/mosquitto.conf -v &
	@cd gateway && make -s start
	@cd admin_panel && make dev

stop-dev :
	-(kill -INT $$(cat pid.tmp)) && rm pid.tmp

.PHONY : build-front
.PHONY : clean-front
.PHONY : run-front
.PHONY : stop-front
.PHONY : run-dev
.PHONY : stop-dev