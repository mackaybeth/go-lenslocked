# We want to build tailwind for production
FROM node:latest AS tailwind-builder
WORKDIR /tailwind
RUN npm init -y && \
    npm install tailwindcss && \
    npx tailwindcss init
# Do not need to watch for changes to the files so no need to mount volumes in docker-compose.yml, just copy
COPY ./templates /templates
COPY ./tailwind/tailwind.config.js /src/tailwind.config.js
COPY ./tailwind/styles.css /src/styles.css
# output is different here from the local dev build
RUN npx tailwindcss -c /src/tailwind.config.js -i /src/styles.css -o /styles.css --minify



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
COPY .env .env
# Copy the binary from the build stage to this stage
COPY --from=builder /app/server ./server
# Copy the built tailwind assets from the build stage to this stage
COPY --from=tailwind-builder /styles.css ./assets/styles.css
CMD ["./server"]
