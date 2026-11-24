FROM golang:1.27.4-alpine3.24 AS builder

ARG BUILD_DATE=$(date "+%Y-%m-%dT%H:%M:%S")

WORKDIR /build

COPY . /build

RUN echo "$( ls -lha )"

RUN go mod tidy && CGO_ENABLED=0 go build \
    -o /build/app \
    -ldflags "-X 'github.com/eldius/mineserver-manager/internal/buildinfo.BuildDate=${BUILD_DATE}' -X 'github.com/eldius/mineserver-manager/internal/buildinfo.Version="dev"  -X github.com/eldius/mineserver-manager/internal/buildinfo.CommitDate=${BUILD_DATE}' -X 'github.com/eldius/mineserver-manager/internal/buildinfo.Commit=dev'" \
        ./cmd/api/

FROM alpine:3.24.4

WORKDIR /app

COPY --from=builder /build/app .
COPY --from=builder /build/config.yaml .

ENTRYPOINT ["/app/app", "-p", "8080"]
