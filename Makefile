build-front :
	(cd gateway; make gateway)
	(cd admin_panel; make build)
	./conf_init.sh
	./ca_gen.sh

run-front :
	@mosquitto -c ./broker/mosquitto.conf -d &> /dev/null
	@cd gateway && make -s start
	@cd admin_panel && make -s start

stop-front :
	-@(kill -INT $$(ps ax | grep mosquitto | grep -v grep | awk '{print $$1}'))
	@cd gateway && make -s stop
	@cd admin_panel && make -s stop

.PHONY : build-front
.PHONY : run-front
.PHONY : stop-front