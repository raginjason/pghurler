FROM golang:1.13.0-alpine as builder

# Install SSL ca certificates.
# ca-certificates is required to call HTTPS endpoints.
# musl-dev is needed for gcc
RUN apk update && apk add --no-cache \
    ca-certificates \
    gcc \
    git \
    musl-dev

WORKDIR /build

# Get dependencies setup first
COPY go.mod go.sum ./
RUN go mod download && go mod verify

ARG APPNAME=pghurler
ENV APPNAME="${APPNAME}"

# Copy rest of application in place and build
COPY . .
RUN GIT_COMMIT=$(git rev-list -1 HEAD) && \
    GIT_ORIGIN=$(git remote get-url origin) && \
    NOW=$(date +'%Y-%m-%d_%T') && \
    CGO_ENABLED=0 GOOS=linux go build -o pghurler -ldflags "-X github.com/raginjason/pghurler/cmd.gitVersion=$GIT_COMMIT -X github.com/raginjason/pghurler/cmd.buildTime=$NOW -X github.com/raginjason/pghurler/cmd.gitOrigin=$GIT_ORIGIN" main.go
# "go test -v all" and "go mod tidy" is suggested to run before release
RUN go test -v ./... -cover

# Bare minimum container
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /app
COPY --from=builder /build/pghurler .

ENTRYPOINT [ "/app/pghurler" ]
