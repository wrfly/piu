FROM golang:alpine
COPY . /go/src/piu
RUN cd /go/src/piu && make build