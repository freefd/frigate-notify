FROM golang:1.22-alpine as build

WORKDIR /app

COPY go.mod /app/
COPY go.sum /app/

RUN go mod download

COPY . /app/

RUN go build -o /frigate-notify .

FROM scratch

WORKDIR /app

COPY --from=build /frigate-notify /app/frigate-notify
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT [ "/app/frigate-notify" ]
