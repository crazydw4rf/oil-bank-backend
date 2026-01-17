FROM golang:latest AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY ./internal ./internal
COPY ./cmd ./cmd
RUN go build -o oil-bank -ldflags="-s -w" ./cmd/oil-bank

COPY db ./db
RUN go build -o dbman -ldflags="-s -w" ./db/...

FROM golang:latest AS runner

WORKDIR /app

COPY --from=builder /build/oil-bank .
COPY --from=builder /build/dbman .
COPY db ./db

CMD ["./oil-bank"]
