FROM alpine:latest

RUN apk add go chromium

WORKDIR /go/src/gxss
COPY . .

RUN go get -d -v ./...
RUN go build .

ENTRYPOINT ["./gxss"]
