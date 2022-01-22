# Build container
FROM golang:1.17-bullseye AS build

# Set build workdir
WORKDIR /app

# Copy app sources
COPY . .

# Build app
RUN go build -o app .

# ---
# Production container
FROM debian:bullseye-slim

# Set app workdir
WORKDIR /app

# Copy binary
COPY --from=build /app/app .

# Run app
CMD ["./app"]
