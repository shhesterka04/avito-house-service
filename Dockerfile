FROM alpine:latest

WORKDIR /root/

# Copy the pre-built binary from the local ./bin/ directory
COPY ./bin/house-service /bin/house-service


# Run the binary
CMD ["./house-service"]