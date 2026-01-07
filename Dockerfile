FROM golang:alpine AS build

RUN mkdir -p /opt/app

WORKDIR /opt/app

COPY . .

RUN apk add build-base

RUN go mod download && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags musl -ldflags="-w -s" \
    -o bookmark_service cmd/api/main.go

FROM alpine AS run

WORKDIR /app

COPY --from=build /opt/app/bookmark_service .

COPY --from=build /opt/app/docs .


CMD ["/app/bookmark_service"]
 