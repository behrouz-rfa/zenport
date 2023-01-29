FROM golang:1.19.4-alpine3.17 AS builder

WORKDIR /zenports

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY .. ./

# Build the binary.
RUN go build -ldflags="-s -w" -v -o monolith ./cmd/zenports

FROM alpine:latest AS runtime

COPY --from=builder /zenports/docker/wait-for .
RUN chmod +x /wait-for

# Copy the binary to the production image from the builder stage.
COPY --from=builder /zenports/monolith /zenports/monolith

CMD ["/zenports/monolith"]
