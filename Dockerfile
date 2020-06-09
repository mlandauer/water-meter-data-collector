FROM golang:1.14 as builder

WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go .

RUN CGO_ENABLED=0 go install -v ./...

FROM scratch

# COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/water-meter-data-collector /go/bin/water-meter-data-collector
# COPY --from=builder --chown=yinyo:0 /tmp /tmp

CMD ["/go/bin/water-meter-data-collector"]