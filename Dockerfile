## Compile the Go binary
FROM golang:1.24.4-alpine AS build-go

# Install git
RUN apk add --no-cache git

ARG GITHUB_USER
ARG GITHUB_PASS

WORKDIR /src

COPY go.mod .
COPY go.sum .
ENV GOPROXY=https://proxy.golang.org
ENV GOPRIVATE="github.com/invopop"

# Prepare credentials for Go repos
RUN git config --global url."https://${GITHUB_USER}:${GITHUB_PASS}@github.com".insteadOf "https://github.com"

RUN go mod download

ADD . /src
RUN go build -o popapp .

## Build Final Container

FROM alpine
RUN apk add --update --no-cache ca-certificates tzdata
WORKDIR /app

COPY --from=build-go /src/popapp /app/
COPY config/config.yaml /app/config/

VOLUME ["/app/config"]

EXPOSE 8080

ENTRYPOINT [ "./popapp" ]
CMD [ "serve" ]
