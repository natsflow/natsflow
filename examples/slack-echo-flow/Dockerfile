#Img build
FROM golang:1.11 AS builder
WORKDIR /slack-echo-flow
COPY . ./
RUN go test -v -mod=readonly ./...
RUN CGO_ENABLED=0 go build -mod=readonly

# Compile statically and use scratch (blank) image to reduce img size
FROM alpine:latest
WORKDIR /slack-echo-flow
COPY --from=builder /slack-echo-flow/slack-echo-flow .
CMD ["./slack-echo-flow"]
