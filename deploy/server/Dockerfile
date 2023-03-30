FROM golang:1.20-alpine AS build
WORKDIR /app

COPY . /app
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s' -tags=go_json,nomsgpack -buildvcs=false -o app ./cmd/auth


FROM alpine:3.17 AS runtime
WORKDIR /app

COPY --from=build /app/web /app/web
COPY --from=build /app/app /usr/bin/

EXPOSE 8080
ENV PORT 8080
CMD ["app"]
