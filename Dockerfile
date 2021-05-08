# Собираем в гошке
FROM golang:1.16.2 as build

ENV BIN_FILE /opt/sysmon/sysmon
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/sysmon/*

# На выходе тонкий образ
FROM ubuntu:21.04

LABEL SERVICE="sysmon"

ENV BIN_FILE "/opt/sysmon/sysmon"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/sysmon/config.toml
COPY ./configs/config.toml ${CONFIG_FILE}

EXPOSE 8080

ENTRYPOINT ${BIN_FILE} -port 8080 -config ${CONFIG_FILE}
