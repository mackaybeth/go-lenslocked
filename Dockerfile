FROM golang as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
# Copy the source code into the container's working directory (the /app dir)
# Do this after setting up the dependencies, as they may not changes as often as the source code
COPY . .
RUN go build -v -o ./server ./cmd/server/

# Start a new container to do these commands
FROM ubuntu
WORKDIR / 
COPY ./assets ./assets
COPY .env .env

# Copy the binary from the first stage to this stage
COPY --from=builder /app/server ./server
CMD ["./server"]
