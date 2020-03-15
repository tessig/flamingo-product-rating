# Example Dockerfile for Flamingo/Go based Projects

# Builder
FROM golang:alpine AS builder
RUN apk update && apk add --no-cache ca-certificates tzdata git && update-ca-certificates
COPY . /app
RUN cd /app && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o rating .

# Final image
FROM scratch

# add timezone data and ssl root certificates
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# add artifacts
ADD VERSION .
ADD config /config
ADD static /static
ADD templates /templates
ADD sql /sql

# add binary
COPY --from=builder /app/rating /rating

ENTRYPOINT ["/rating"]
CMD ["serve"]
