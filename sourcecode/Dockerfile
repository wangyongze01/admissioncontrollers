#FROM newharbor.ipinyou.com/deepzero/openjdk:8u342-jdk-slimaa
#
#COPY webhook /usr/local/bin/webhook
#RUN chmod +x /usr/local/bin/webhook
#
#ENTRYPOINT ["webhook"]

#golang:1.22.1-alpine
FROM newharbor.ipinyou.com/alphadata/golang:1.22.1-alpine as dev-env

WORKDIR /app

FROM dev-env as build-env
COPY go.mod /go.sum /app/
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go env -w GOPROXY=https://goproxy.cn,direct

RUN go mod download


COPY . /app/

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 go build -o /webhook

FROM newharbor.ipinyou.com/alphadata/alpine:3.20.2 as runtime

COPY --from=build-env /webhook /usr/local/bin/webhook
RUN chmod +x /usr/local/bin/webhook

ENTRYPOINT ["webhook"]