FROM golang:1.20 AS builder

WORKDIR /go/src/app
COPY . .

RUN apt update && apt install -y upx
RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -o metrics
RUN upx metrics

FROM debian:10 AS runner

WORKDIR /go/bin
EXPOSE 8855
COPY --from=builder /go/src/app/metrics /go/bin/metrics
HEALTHCHECK --timeout=10s CMD curl --silent --fail http://127.0.0.1:8088/healthcheck

CMD ["./metrics", "-p", "8855", "-interval", "30", "-mcAddress", "172.124.30.1", "-mcPort", "25565"]