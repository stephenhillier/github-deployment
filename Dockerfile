FROM golang:1.11-stretch AS builder
WORKDIR /go/src/github.com/stephenhillier/github-deployment
ADD . .
RUN go test && go build

FROM debian:stretch
RUN mkdir /app
COPY --from=builder /go/src/github.com/stephenhillier/github-deployment/github-deployment /app/
ENTRYPOINT [ "/app/github-deployment" ]
