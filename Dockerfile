# build stage
FROM golang:1.16-alpine AS build-env
ADD . /src
# ginserver
RUN cd /src/ginserver && go mod download
RUN cd /src/ginserver && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ginserver
# httpserver
RUN cd /src/httpserver && go mod download
RUN cd /src/httpserver && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o httpserver
# grpcserver


# final stage
FROM alpine
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-env /src/ginserver/ginserver /ginserver
COPY --from=build-env /src/httpserver/httpserver /httpserver
# COPY --from=build-env /src/grpcserver/grpcserver /grpcserver