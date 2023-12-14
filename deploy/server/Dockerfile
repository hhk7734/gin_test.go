FROM golang:1.20-alpine AS build
WORKDIR /app

COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -a \
    -ldflags '-s' \
    -tags=go_json,nomsgpack \
    -buildvcs=false \
    -o app \
    ./cmd/server


FROM alpine:3.17 AS runtime
WORKDIR /app

COPY --from=build /app/web /app/web
COPY --from=build /app/app /usr/bin/

ENTRYPOINT [ "app" ]