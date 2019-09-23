FROM wrfly/piu:build-base
COPY . /src
RUN cd /src && \
    make build && \
    mv bin/ /artifacts

FROM alpine
ARG BINARY
RUN apk add --update ca-certificates
COPY --from=0 /artifacts/${BINARY} /usr/local/bin/
CMD piu --debug
