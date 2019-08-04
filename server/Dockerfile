FROM golang:1.12-alpine3.10 as builder

RUN apk --update add alpine-sdk

WORKDIR /src

ADD Makefile ./
RUN make data

ADD go.mod go.sum ./
RUN go mod download

ADD . .
RUN make build

FROM alpine:3.10 as app

RUN apk update && \
  apk --update add ca-certificates && \
  update-ca-certificates && \
  apk add curl cyrus-sasl-dev openssl-dev

COPY --from=builder /src/bin/adsb-server /usr/bin
COPY --from=builder /src/data/StandingData.sqb /opt/StandingData.sqb

EXPOSE 8080

CMD ["adsb-server", "--addr", ":8080", "--db", "/opt/StandingData.sqb"]
