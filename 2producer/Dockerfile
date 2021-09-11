FROM golang:alpine as build
WORKDIR /go/src/app
# Set necessary environmet variables needed for our image
ENV GOOS=linux GOARCH=amd64 
COPY ./app .
RUN apk update && apk add gcc librdkafka-dev openssl-libs-static zlib-static zstd-libs libsasl librdkafka-static lz4-dev lz4-static zstd-static libc-dev musl-dev 
# RUN go build -tags musl -ldflags '-w -extldflags "-static"' -o main

# FROM scratch
# COPY --from=build /go/src/app/main /
# # Command to run the executable
# ENTRYPOINT ["/main"]