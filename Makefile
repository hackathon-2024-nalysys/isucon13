MAKE=make -C

DOCKER_BUILD=docker build
DOCKER_BUILD_OPTS=--no-cache
DOCKER_RMI=docker rmi -f

ISUPIPE_TAG=isupipe:latest

SERVER_IP := 18.181.199.55

test: test_benchmarker
.PHONY: test bench

test_benchmarker:
	$(MAKE) bench test
.PHONY: test_benchmarker

build_webapp:
	$(MAKE) webapp/go docker_image
.PHONY: build_webapp

bench:
	ssh isubench ./bench run --target https://pipe.u.isucon.local --nameserver $(SERVER_IP) --webapp $(SERVER_IP) --enable-ssl

deploy_benchmarker:
	scp ./Makefile "isubench:~/Makefile"