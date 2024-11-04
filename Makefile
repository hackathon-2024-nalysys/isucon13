MAKE=make -C

DOCKER_BUILD=docker build
DOCKER_BUILD_OPTS=--no-cache
DOCKER_RMI=docker rmi -f

ISUPIPE_TAG=isupipe:latest

SERVER_IP := 18.181.199.55

SOURCE_DIR := ~/webapp/go

test: test_benchmarker
.PHONY: test bench

test_benchmarker:
	$(MAKE) bench test
.PHONY: test_benchmarker

build_webapp:
	$(MAKE) webapp/go docker_image
.PHONY: build_webapp alp query pprof

bench:
	make reset
	ssh isubench ./bench run --target https://pipe.u.isucon.local --nameserver $(SERVER_IP) --webapp $(SERVER_IP) --enable-ssl

reset:
	ssh isucon1 "sudo killall -USR2 isupipe && sudo rm /var/log/nginx/access.log && sudo systemctl restart nginx && sudo rm -f /var/log/mysql/mysql-slow.log && sudo systemctl restart mysql"

deploy_benchmarker:
	scp ./Makefile "isubench:~/Makefile"

deploy:
	$(MAKE) webapp/go build
	ssh isucon1 rm -rf "$(SOURCE_DIR)"
	scp -r webapp/go isucon1:"$(SOURCE_DIR)"
	scp ./Makefile "isucon1:~/Makefile"
	ssh isucon1 sudo systemctl restart isupipe-go

alp:
	sudo cat /var/log/nginx/access.log | alp ltsv --sort sum -r -m "user/\w+/icon,user/\w+/statistics,user/\w+/theme,assets/.*,livestream/\d+/report,livestream/\d+/livecomment,livestream/\d+/reaction,livestream/\d+/ngwords,livestream/\d+/enter,livestream/\d+/exit,livestream/\d+/moderate,livestream/\d+/moderate,livestream/\d+/statistics,livestream/\d+$$" -o count,method,uri,min,avg,max,sum | less

query:
	sudo pt-query-digest /var/log/mysql/mysql-slow.log

pprof:
	killall -USR1 isupipe && \
	sleep 1 && \
	$(HOME)/go/bin/pprof -http=localhost:1080 "$(SOURCE_DIR)"/isupipe "$(SOURCE_DIR)"/cpu.pprof