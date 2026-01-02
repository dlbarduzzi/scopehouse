# syntax=docker/dockerfile:1

ARG GO_VERSION=1.25.5

# Initial stage.
FROM public.ecr.aws/docker/library/golang:${GO_VERSION} AS base
LABEL org.opencontainers.image.source=https://github.com/dlbarduzzi/scopehouse

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

# Builder stage.
FROM base AS builder
LABEL org.opencontainers.image.source=https://github.com/dlbarduzzi/scopehouse

COPY . /app

RUN CGO_ENABLED=0 go build -o /bin/scopehouse ./cmd/scopehouse

# Running stage.
FROM gcr.io/distroless/static-debian12
LABEL org.opencontainers.image.source=https://github.com/dlbarduzzi/scopehouse

USER nonroot:nonroot
COPY --from=builder --chown=nonroot:nonroot /bin/scopehouse /bin/scopehouse

ENTRYPOINT ["/bin/scopehouse"]
