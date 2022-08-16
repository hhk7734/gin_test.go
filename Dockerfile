FROM golang:1.18.1-alpine AS build
COPY . /app
WORKDIR /app
RUN apk --update add ca-certificates
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s' -tags=go_json,nomsgpack -buildvcs=false -o app

FROM scratch AS runtime
WORKDIR /app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /app/app /usr/bin/
CMD ["app"]