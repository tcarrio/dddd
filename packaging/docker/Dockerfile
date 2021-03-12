FROM docker.io/library/golang:latest as build_image

COPY . /app/
WORKDIR /app/
RUN make micro

FROM scratch

COPY --from=build_image /app/release/microdddd /dddd
COPY --from=build_image /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/dddd", "-c", "/opt/dddd/config.yml"]
