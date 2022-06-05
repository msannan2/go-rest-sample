ARG BUILDER_BASE=golang:1.16-alpine3.12
FROM ${BUILDER_BASE} AS builder

WORKDIR /app

RUN apk add --update --no-cache make git

COPY go.mod .

RUN go mod download

COPY . .

RUN make cmd/sample/sample

FROM alpine:3.12

RUN apk add --update --no-cache dumb-init

ENV APP_PORT=59000
EXPOSE $APP_PORT

WORKDIR /
COPY --from=builder /app/cmd/sample/sample /
COPY --from=builder /app/cmd/sample/res/configuration.toml /res/configuration.toml

ENTRYPOINT ["/sample"]