FROM golang:1.15-alpine AS build
ARG GH_CI_TOKEN=$GH_CI_TOKEN

RUN apk add git bash pkgconfig vips-dev vips gcc musl-dev

WORKDIR /app
COPY / /app
ENV GOPRIVATE="github.com/nnqq/*"
RUN git config --global url."https://nnqq:$GH_CI_TOKEN@github.com/".insteadOf "https://github.com/"
RUN CGO_CFLAGS_ALLOW=-Xpreprocessor go build -o servicebin

RUN GRPC_HEALTH_PROBE_VERSION=v0.3.2 && \
    wget -qO/app/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /app/grpc_health_probe

FROM alpine:latest

RUN apk --update upgrade && \
    apk add curl ca-certificates && \
    apk add vips-dev vips && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/*

WORKDIR /app
COPY --from=build /app/servicebin /app
COPY --from=build /app/grpc_health_probe /app
