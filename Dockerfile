# Dockerfile References: https://docs.docker.com/engine/reference/builder/

ARG GO_VERSION=1.12.7-alpine3.9

FROM golang:${GO_VERSION} as builder

# Git is required for fetching the dependencies.
RUN apk add --no-cache ca-certificates git

# Add Maintainer Info
LABEL maintainer="Peshiy Sergey <pesheysergey@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /go/src/github.com/speshiy/V-K-Alcohol-Excise-Parser

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/vkaep_server .

######## Start a new stage from scratch #######
FROM alpine:3.10 as final 

ENV RELEASE=$RELEASE
ENV SSL=$SSL
ENV SCT=$SCT
ENV WINDBHOST=$WINDBHOST
ENV TZ=Europe/London

RUN apk add --no-cache tzdata

WORKDIR /home/tuvis-server

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /go/bin/vkaep_server .

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /go/src/github.com/speshiy/Tuvis-Server/frontend/ ./frontend
# Import the root ca-certificates (required for Let's Encrypt)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 80
EXPOSE 443
EXPOSE 8081

# Mount the certificate cache directory as a volume, so it remains even after
# we deploy a new version
VOLUME ["cert-cache"]

ENTRYPOINT ./vkaep_server -release=$RELEASE -ssl=$SSL -sct=$SCT -windbhost=$WINDBHOST