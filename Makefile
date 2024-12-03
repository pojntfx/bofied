# Public variables
DESTDIR ?=
PREFIX ?= /usr/local
OUTPUT_DIR ?= out
DST ?=

WWWROOT ?= /var/www/html
WWWPREFIX ?= /bofied

# Private variables
clis = bofied-backend
pwas = frontend
all: build

# Build
build: build/cli build/pwa

build/cli: $(addprefix build/cli/,$(clis))
$(addprefix build/cli/,$(clis)):
ifdef DST
	go build -o $(DST) ./cmd/$(subst build/cli/,,$@)
else
	go build -o $(OUTPUT_DIR)/$(subst build/cli/,,$@) ./cmd/$(subst build/cli/,,$@)
endif

build/pwa:
	GOOS=js GOARCH=wasm go build -o web/app.wasm cmd/bofied-frontend/main.go
	rm -rf $(OUTPUT_DIR)/bofied-frontend/
	mkdir -p $(OUTPUT_DIR)/bofied-frontend/web
	go run ./cmd/bofied-frontend/main.go --build
	cp -r web/* $(OUTPUT_DIR)/bofied-frontend/web
	tar -cvzf $(OUTPUT_DIR)/frontend.tar.gz -C $(OUTPUT_DIR)/bofied-frontend .

# Special target for GitHub pages builds
build/pwa-github-pages:
	GOOS=js GOARCH=wasm go build -o web/app.wasm cmd/bofied-frontend/main.go
	rm -rf $(OUTPUT_DIR)/bofied-frontend/
	mkdir -p $(OUTPUT_DIR)/bofied-frontend/web
	go run ./cmd/bofied-frontend/main.go --build --path bofied
	cp -r web/* $(OUTPUT_DIR)/bofied-frontend/web
	tar -cvzf $(OUTPUT_DIR)/frontend.tar.gz -C $(OUTPUT_DIR)/bofied-frontend .

# Install
install: install/cli install/pwa

install/cli: $(addprefix install/cli/,$(clis))
$(addprefix install/cli/,$(clis)):
	install -D -m 0755 $(OUTPUT_DIR)/$(subst install/cli/,,$@) $(DESTDIR)$(PREFIX)/bin/$(subst install/cli/,,$@)

install/pwa:
	mkdir -p $(DESTDIR)$(WWWROOT)$(WWWPREFIX)
	tar -xvf out/frontend.tar.gz -C $(DESTDIR)$(WWWROOT)$(WWWPREFIX)

# Uninstall
uninstall: uninstall/cli uninstall/pwa

uninstall/cli: $(addprefix uninstall/cli/,$(clis))
$(addprefix uninstall/cli/,$(clis)):
	rm $(DESTDIR)$(PREFIX)/bin/$(subst uninstall/cli/,,$@)

uninstall/pwa:
	rm -rf $(DESTDIR)$(WWWROOT)$(WWWPREFIX)

# Run
run: run/cli run/pwa

run/cli: $(addprefix run/cli/,$(clis))
$(addprefix run/cli/,$(clis)):
	$(subst run/cli/,,$@) $(ARGS)

run/pwa:
	go run ./cmd/bofied-frontend/ --serve

# Dev
dev: dev/cli dev/pwa

dev/cli: $(addprefix dev/cli/,$(clis))
$(addprefix dev/cli/,$(clis)): $(addprefix build/cli/,$(clis))
	sudo setcap cap_net_bind_service+ep $(OUTPUT_DIR)/$(subst dev/cli/,,$@)
	$(OUTPUT_DIR)/$(subst dev/cli/,,$@) $(ARGS)

dev/pwa: build/pwa
	go run ./cmd/bofied-frontend/ --serve

# Test
test: test/cli test/pwa

test/cli:
	go test -timeout 3600s -parallel $(shell nproc) ./...

test/pwa:
	true

# Benchmark
benchmark: benchmark/cli benchmark/pwa

benchmark/cli:
	go test -timeout 3600s -bench=./... ./...

benchmark/pwa:
	true

# Clean
clean: clean/cli clean/pwa

clean/cli:
	rm -rf out pkg/models pkg/api/proto/v1 rm -rf ~/.local/share/bofied

clean/pwa:
	rm -rf out web/app.wasm

# Dependencies
depend: depend/cli depend/pwa

depend/cli:
	GO111MODULE=on go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	GO111MODULE=on go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	GO111MODULE=on go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

	go generate ./...

depend/pwa:
	true
