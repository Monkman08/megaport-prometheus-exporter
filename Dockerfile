FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o megaport-prometheus-exporter

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /app/megaport-prometheus-exporter /megaport-prometheus-exporter

ENTRYPOINT ["/megaport-prometheus-exporter"]