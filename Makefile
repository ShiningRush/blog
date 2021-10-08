SRV_NAME=blog
HUB_URL=shiningrush

.PHONY: init
init:
	git submodule update --init --recursive

.PHONY: build
build:
	hugo
	docker build . -t $(HUB_URL)/$(SRV_NAME)

.PHONY: publish
publish: build
	docker push $(HUB_URL)/$(SRV_NAME)