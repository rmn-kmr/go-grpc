NAME := lsp
DOCKER_REPO := gcr.io/rmnkmr-42/$(NAME)

###################
# TARGETS
#################
default: run

run:
	@go run cmd/main.go

job:
	@go run cmd/main.go job $(JOB)

build: OUT ?= $(NAME)
build:
	@CGO_ENABLED=0 go build -o bin/$(OUT) cmd/main.go

docker:
	docker build --build-arg SSH_PRIVATE_KEY="$$(cat $(HOME)/.ssh/id_rsa)" -t $(NAME) -f Dockerfile ./..

docker-run: docker
docker-run:
	docker run \
		-v "$$(pwd)/config.yml:/config.yml" \
		-v "$$(pwd)/service-account.json:/service-account.json" \
		-p "8080:8080" \
		$(NAME) server

docker-push: TAG ?= latest
docker-push: docker
	docker tag $(NAME) $(DOCKER_REPO):$(TAG)
	docker push $(DOCKER_REPO):$(TAG)


.PHONY: tools
tools: # Install dependencies and tools required to build
	@echo "Fetching tools..."
	@go generate -tags tools tools/tools.go
	@echo
	@echo "Done!"

.PHONY: gen/server
gen/server: # Generates grpc protobuf Go files from grpc.proto
	@go generate ./pkg/server/server.go

.PHONY: gen/client
gen/client: # Generates grpc-gateway client from grpc.proto
	@mkdir -p pkg/client/go
	@go generate ./pkg/client/client.go


.PHONY: gen/doc
gen/doc: # generates the grpc proto docs
	@go generate ./docs

.PHONY: proto
proto:
	@buf generate

.PHONY: sql
sql:
	@sqlc generate

# Format Go files, ignoring files marked as generated through the header defined at
# https://pkg.go.dev/cmd/go#hdr-Generate_Go_files_by_processing_source
.PHONY: fmt
fmt:
	@grep -L -R "^\/\/ Code generated .* DO NOT EDIT\.$$" --exclude-dir=.git --include="*.go" . | xargs gofmt -w
	@buf format -w
