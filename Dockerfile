FROM golang
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
# Copy the source code into the container's working directory (the /app dir)
# Do this after setting up the dependencies, as they may not changes as often as the source code
COPY . .
RUN go build -v -o ./server ./cmd/server/
CMD ./server
