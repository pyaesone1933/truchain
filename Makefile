PACKAGES=$(shell go list ./...)

MODULES = argument backing category challenge expiration stake story

CHAIN_DIR = ./.chain

define \n


endef

benchmark:
	@go test -bench=. $(PACKAGES)

buidl: build

build: build_daemon

build_cli:
	go build -o bin/trucli cmd/trucli/main.go

build_daemon:
	go build -o bin/truchaind cmd/truchaind/main.go

build_linux:
	GOOS=linux GOARCH=amd64 go build -o bin/truchaind cmd/truchaind/main.go

doc:
	@echo "--> Wait a few seconds and visit http://localhost:6060/pkg/github.com/TruStory/truchain/"
	godoc -http=:6060

init:
	bin/truchaind --home $(CHAIN_DIR) init

reset:
	bin/truchaind --home $(CHAIN_DIR) unsafe-reset-all

restart: build_daemon reset start

start:
	bin/truchaind --home $(CHAIN_DIR) --log_level "main:info,state:info,*:error,app:info,backing:info,category:info,challenge:info,expiration:info,stake:info,stories:info,vote:info,voting:info" start

check:
	gometalinter ./...

dep_graph: ; $(foreach dir, $(MODULES), godepgraph -s -novendor github.com/TruStory/truchain/x/$(dir) | dot -Tpng -o x/$(dir)/dep.png${\n})

install_tools_macos:
	brew install dep && brew upgrade dep
	brew tap alecthomas/homebrew-tap && brew install gometalinter

go_test:
	@go test $(PACKAGES)

set_registrar:
	mkdir -p $(HOME)/.truchaind
	cp $(CHAIN_DIR)/registrar.key $(HOME)/.truchaind/registrar.key

set_env_vars:
	mkdir -p $(HOME)/.truchaind
	cp .env.example $(CHAIN_DIR)/.env
	cp .env.example $(HOME)/.truchaind/.env

test: set_registrar set_env_vars go_test

test_cover: set_registrar set_env_vars
	@go test $(PACKAGES) -v -timeout 30m -race -coverprofile=coverage.txt -covermode=atomic
	@go tool cover -html=coverage.txt

update_deps:
	@echo "--> Running dep ensure"
	@rm -rf .vendor-new
	@dep ensure -v -vendor-only

.PHONY: benchmark buidl build check dep_graph test test_cover update_deps

