FROM golang:1.11.1-alpine as builder

ENV GOPATH /go

COPY . /go/src/github.com/spark-cluster

# Build
RUN  go build -o /spark-cluster-operator /go/src/github.com/spark-cluster/cmd/manager && \
	 go build -o /backend /go/src/github.com/spark-cluster/dashboard/backend

FROM alpine:3.8

RUN apk upgrade --update --no-cache

COPY --from=builder /spark-cluster-operator /usr/local/bin/spark-cluster-operator
COPY --from=builder /backend /usr/local/bin/backend
COPY --from=builder /go/src/github.com/spark-cluster/dashboard/frontend /usr/local/frontend
