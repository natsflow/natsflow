#Img build
FROM golang:1.11 AS builder
WORKDIR /echo-nats
COPY . ./
RUN go test -v -mod=readonly ./...
RUN CGO_ENABLED=0 go build -mod=readonly
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -a -installsuffix cgo -ldflags '-extldflags "-static"'

# Compile statically and use scratch (blank) image to reduce img size
FROM scratch
WORKDIR /echo-nats
COPY --from=builder /echo-nats/echo-nats .
CMD ["./echo-nats"]
