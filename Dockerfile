FROM golang:1.11-alpine3.9
LABEL "name"="github-deployment"
LABEL "version"="0.1.0"

RUN mkdir -p /go/src/github-deployment/
COPY . /go/src/github-deployment

RUN go install -v /go/src/github-deployment/

ENTRYPOINT ["/go/src/github-deployment/entrypoint.sh"]
