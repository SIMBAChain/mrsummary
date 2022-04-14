# syntax=docker/dockerfile:1
FROM golang:alpine as app-builder
WORKDIR /go/src/github.com/simbachain/mrsummary/
RUN apk add git gcc musl-dev
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go install -ldflags '-extldflags "-static"' -tags timetzdata
RUN chmod +x /go/bin/MRSummary

FROM busybox:latest
COPY --from=app-builder /go/bin/MRSummary /bin/MRSummary
COPY --from=app-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/bin/MRSummary"]