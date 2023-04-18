PROJDIR=$(dir $(realpath $(firstword $(MAKEFILE_LIST))))

# change to project dir so we can express all as relative paths
$(shell cd $(PROJDIR))

REPO_PATH=github.com/sorintlab/stolon

# VERSION ?= $(shell scripts/git-version.sh)
# LD_FLAGS="-w -X $(REPO_PATH)/cmd.version=$(VERSION)"
LD_FLAGS += -w
LD_FLAGS += -X "$(REPO_PATH)/cmd.version=$(shell git describe --tags --dirty --always)"
LD_FLAGS += -X "$(REPO_PATH)/cmd.commit=$(shell git rev-parse HEAD)"
LD_FLAGS += -X "$(REPO_PATH)/cmd.date=$(shell date -u '+%Y-%m-%d %I:%M:%S')"

$(shell mkdir -p bin )


.PHONY: all
all: build

.PHONY: build
build: sentinel keeper proxy stolonctl

.PHONY: test
test: build
	./test

.PHONY: sentinel keeper proxy stolonctl docker

keeper:
	GO111MODULE=on go build -ldflags '$(LD_FLAGS)' -o $(PROJDIR)/bin/stolon-keeper $(REPO_PATH)/cmd/keeper

sentinel:
	CGO_ENABLED=0 GO111MODULE=on go build -ldflags '$(LD_FLAGS)' -o $(PROJDIR)/bin/stolon-sentinel $(REPO_PATH)/cmd/sentinel

proxy:
	CGO_ENABLED=0 GO111MODULE=on go build -ldflags '$(LD_FLAGS)' -o $(PROJDIR)/bin/stolon-proxy $(REPO_PATH)/cmd/proxy

stolonctl:
	CGO_ENABLED=0 GO111MODULE=on go build -ldflags '$(LD_FLAGS)' -o $(PROJDIR)/bin/stolonctl $(REPO_PATH)/cmd/stolonctl

.PHONY: docker
docker:
	if [ -z $${PGVERSION} ]; then echo 'PGVERSION is undefined'; exit 1; fi; \
	if [ -z $${TAG} ]; then echo 'TAG is undefined'; exit 1; fi; \
	docker build --build-arg PGVERSION=${PGVERSION} -t ${TAG} -f examples/kubernetes/image/docker/Dockerfile .
