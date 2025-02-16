# Build Golang API
FROM golang:1.23-alpine AS gobuilder

WORKDIR /app
RUN apk --no-cache add ca-certificates make
COPY go.mod .
COPY go.sum .
COPY . .

RUN make install

RUN make build


# Build Angular Frontend
FROM node:23.0-alpine AS nodebuilder

WORKDIR /app
COPY web .

RUN npm cache clean --force
RUN npm install -g @angular/cli
RUN npm install
RUN npm run build

# Final Build
FROM alpine:latest

WORKDIR /app

# Copy Golang binary
COPY --from=gobuilder /app/tmp/main /app/main

# Copy Angular dist files
COPY --from=nodebuilder /app/dist/web/browser /app/web/dist

# Expose necessary ports
EXPOSE 3000

# Command to run the Golang API
ENTRYPOINT ["/app/main"]