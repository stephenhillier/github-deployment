FROM golang:1.11-stretch AS builder
WORKDIR /go/src/github.com/stephenhillier/github-deployment
ADD . .
RUN go test && go build

FROM alpine:3.9
RUN mkdir /app
COPY --from=builder /go/src/github.com/stephenhillier/github-deployment/github-deployment /app/
RUN chmod +x /app/api
ENTRYPOINT [ "/app/github-deployment" ]
