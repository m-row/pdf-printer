include .env
export

# Go
.PHONY: list update init
list:
	go list -m -u
update:
	go get -u ./...
init:
	@go install github.com/cosmtrek/air@latest
	@go install github.com/go-delve/delve/cmd/dlv@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.1
	@go install github.com/jesseduffield/lazygit@latest
	@go install github.com/nametake/golangci-lint-langserver@latest
	@go install github.com/segmentio/golines@latest
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@go install mvdan.cc/gofumpt@latest

CONTAINER_NAME:=${CONTAINER_NAME}
CONTAINER_TAG:=${APP_VER}.$(shell git rev-list --count HEAD).$(shell git describe --always)
CONTAINER_IMG:=${CONTAINER_NAME}:${CONTAINER_TAG}

# hub.docker.com: Production
.PHONY: dh dh/down dh/push prune ps inspect
dh:
	export CONTAINER_TAG=${CONTAINER_TAG}
	docker compose -f docker-compose.yml up --build -d
dh/down:
	docker compose -f docker-compose.yml down
dh/push: dh
	docker tag sadeem/${CONTAINER_IMG} ${CONTAINER_REG}/${CONTAINER_IMG}
	docker tag sadeem/${CONTAINER_IMG} ${CONTAINER_REG}/${CONTAINER_NAME}:latest
	docker push ${CONTAINER_REG}/${CONTAINER_NAME} -a

prune:
	docker system prune -a -f --volumes
ps:
	docker ps --format "table {{.Names}}\t{{.Status}}\t{{.RunningFor}}\t{{.Size}}\t{{.Ports}}"
inspect: 
	docker inspect -f "{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}" $(n)

