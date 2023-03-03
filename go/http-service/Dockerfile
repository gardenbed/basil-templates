# BUILD STAGE
FROM golang:1.20-alpine as builder
RUN apk add --no-cache git
WORKDIR /project
COPY . .
ARG ldflags
RUN go build -o service -ldflags "$ldflags"

# FINAL STAGE
FROM alpine:3.17
EXPOSE 8080
RUN apk add --no-cache curl ca-certificates
HEALTHCHECK --interval=5m --timeout=3s CMD curl -f http://localhost:8080/health || exit 1
COPY --from=builder /project/service /usr/local/bin/
RUN chown -R nobody:nogroup /usr/local/bin/service
USER nobody
ENTRYPOINT [ "service" ]
