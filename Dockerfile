FROM golang:1.22.0 as builder

WORKDIR /app

# Fetch the go mod deps first so docker build can cache these
# making the builds faster.
COPY go.mod go.sum ./
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

RUN go build -o /api ./cmd/api/.

# Start a new stage from scratch:
FROM alpine:3.18

COPY --from=builder /api /api

CMD /api
