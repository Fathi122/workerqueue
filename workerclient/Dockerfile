FROM golang:alpine AS builder

RUN mkdir -p /go/src/workerqueue/workerclient && \
    mkdir -p /go/src/workerqueue/conf && mkdir -p /go/src/workerqueue/workerproto \
    mkdir -p /go/src/workerqueue/certs && \
    mkdir -p /go/bin/certs && mkdir -p /go/bin/conf &&  mkdir -p /go/bin/workerclient
COPY ./workerclient/main.go /go/src/workerqueue/workerclient
COPY ./certs/server.crt /go/src/workerqueue/certs/
COPY ./certs/server.pem /go/src/workerqueue/certs/
COPY ./certs/ca.crt /go/src/workerqueue/certs/

COPY ./conf/ /go/src/workerqueue/conf/

COPY ./certs/server.crt /go/bin/certs/
COPY ./certs/server.pem /go/bin/certs/
COPY ./certs/ca.crt /go/bin/certs/
COPY ./conf/config.yaml /go/bin/conf/

COPY ./workerproto/ /go/src/workerqueue/workerproto/

COPY ./go.mod /go/src/workerqueue/
COPY ./go.sum /go/src/workerqueue/

# set GOPATH & GOBIN
ENV GOPATH /go
ENV GOBIN /go/bin

# RUN build
RUN cd /go/src/workerqueue/ && GOOS=linux GOARCH=amd64 go build -o clientserver workerclient/main.go && mv clientserver /go/bin/workerclient

# alpine based
FROM alpine

RUN mkdir -p /workerclient && mkdir -p /conf && mkdir -p /certs
COPY --from=builder /go/bin/workerclient /workerclient/.
COPY --from=builder /go/bin/conf /conf/.
COPY --from=builder /go/bin/certs /certs/.

# Install CA certificates
RUN apk update && apk add --no-cache --update ca-certificates tzdata && chmod +x /workerclient/clientserver

WORKDIR /workerclient
RUN ls -lart /workerclient/

EXPOSE 8080

CMD [ "/workerclient/clientserver" ]
