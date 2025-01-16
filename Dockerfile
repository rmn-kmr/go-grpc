################################################################
#  BUILD
################################################################
FROM asia-south1-docker.pkg.dev/rmnkmr-42/common/golang:1.21-alpine3.18 AS builder

COPY . /src
WORKDIR /src
RUN go mod download
RUN make build


################################################################
#  MAIN
################################################################
FROM alpine:3.18
RUN echo "https://dl-cdn.alpinelinux.org/alpine/v3.18/main" >/etc/apk/repositories
RUN echo "https://dl-cdn.alpinelinux.org/alpine/v3.18/community" >>/etc/apk/repositories
RUN apk add --no-cache --update curl ca-certificates && update-ca-certificates
RUN apk add jq
RUN mkdir -p /bin
COPY --from=builder /src/bin/lsp /bin/lsp
HEALTHCHECK CMD curl --fail http://0.0.0.0:8080/ || exit 1
ENTRYPOINT ["/bin/lsp"]
EXPOSE 8080
