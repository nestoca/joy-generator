FROM golang:1.22-alpine AS build

ENV CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

RUN openssh-client ca-certificates && update-ca-certificates 2>/dev/null || true

ENV HOME=/home/golang

WORKDIR /app

RUN adduser -h $HOME -D -u 1000 -G root golang && \
  chown golang:root /app && \
  chmod g=u /app $HOME

USER golang:root

COPY --chown=golang:root go.mod go.sum ./

RUN go mod download

COPY --chown=golang:root cmd ./cmd
COPY --chown=golang:root internal ./internal

RUN go build -v -o joy-generator ./cmd/server

FROM alpine:3.18 AS prod

COPY --from=build /etc/passwd /etc/group  /etc/
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build --chown=golang:root /app/joy-generator /app/

RUN apk add helm

USER golang:root
EXPOSE 8080

WORKDIR /app

ENTRYPOINT ["./joy-generator"]
