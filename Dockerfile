FROM alpine:3.2
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
ADD release/linux/amd64/drone-dispatch /bin/
ENTRYPOINT ["/bin/drone-dispatch"]
