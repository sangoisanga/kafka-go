#####################################
#   STEP 1 build executable binary  #
#####################################
FROM golang:alpine AS builder

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go install -v ./...

#####################################
#   STEP 2 build a small image      #
#####################################
FROM scratch

# Copy our static executable.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/bin/* /

EXPOSE 8080

# Run the hello binary.
ENTRYPOINT ["/producer"]