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

WORKDIR /home/nonroot
RUN apk --no-cache add shadow

RUN useradd -m nonroot
RUN mkdir -p /home/nonroot/storage
ENV GOOGLE_APPLICATION_CREDENTIALS="/google-sa-credentials.json"

# Copy Golang binary
COPY --from=gobuilder /app/tmp/main /home/nonroot/main

# Copy Angular dist files
COPY --from=nodebuilder /app/dist/web/browser /home/nonroot/web/dist

RUN chown -R nonroot:nonroot /home/nonroot/main
RUN chown -R nonroot:nonroot /home/nonroot/web
RUN chown -R nonroot:nonroot /home/nonroot/storage

# User to nonroot
USER nonroot

# Expose necessary ports
EXPOSE 3000

# Command to run the Golang API
ENTRYPOINT ["/home/nonroot/main"]