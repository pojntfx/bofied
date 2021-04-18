all: build

backend:
	go build -o out/bofied-backend/bofied-backend cmd/bofied-backend/main.go

frontend:
	rm -f web/app.wasm
	GOOS=js GOARCH=wasm go build -o web/app.wasm cmd/bofied-frontend/main.go
	go build -o /tmp/bofied-frontend-build cmd/bofied-frontend/main.go
	rm -rf out/bofied-frontend
	/tmp/bofied-frontend-build -build
	cp -r web/* out/bofied-frontend/web

build: backend frontend

release-backend:
	CGO_ENABLED=1 go build -ldflags="-extldflags=-static" -tags netgo -o out/release/bofied-backend/bofied-backend.linux-$$(uname -m) cmd/bofied-backend/main.go

release-frontend: frontend
	rm -rf out/release/bofied-frontend
	mkdir -p out/release/bofied-frontend
	cd out/bofied-frontend && tar -czvf ../release/bofied-frontend/bofied-frontend.tar.gz .

release-frontend-github-pages: frontend
	rm -rf out/release/bofied-frontend-github-pages
	mkdir -p out/release/bofied-frontend-github-pages
	/tmp/bofied-frontend-build -build -path bofied -out out/release/bofied-frontend-github-pages
	cp -r web/* out/release/bofied-frontend-github-pages/web

release: release-backend release-frontend release-frontend-github-pages

install: release-backend
	sudo install out/release/bofied-backend/bofied-backend.linux-$$(uname -m) /usr/local/bin
	sudo setcap cap_net_bind_service+ep /usr/local/bin/bofied-backend
	
dev:
	while [ -z "$$BACKEND_PID" ] || [ -n "$$(inotifywait -q -r -e modify pkg cmd web/*.css)" ]; do\
		$(MAKE);\
		kill -9 $$BACKEND_PID 2>/dev/null 1>&2;\
		kill -9 $$FRONTEND_PID 2>/dev/null 1>&2;\
		wait $$BACKEND_PID $$FRONTEND_PID;\
		sudo setcap cap_net_bind_service+ep out/bofied-backend/bofied-backend;\
		out/bofied-backend/bofied-backend & export BACKEND_PID="$$!";\
		/tmp/bofied-frontend-build -serve & export FRONTEND_PID="$$!";\
	done

clean:
	rm -rf out
	rm -rf ~/.local/share/bofied

depend:
	# Generate bindings
	go generate ./...