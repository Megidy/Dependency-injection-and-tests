FROM golang:1.23-alpine AS builder 

WORKDIR /app

RUN apk --no-cache add bash git make gcc gettext musl-dev

COPY ["go.mod","go.sum","./"]
RUN go mod download

COPY . .

# api
RUN go build -o ./api ./cmd/api/main.go

# depot 
RUN go build -o ./depot ./cmd/depot/main.go

# tests
FROM builder AS run-test-stage
    RUN go test -v ./...

FROM alpine

WORKDIR /
# api build
COPY --from=builder /app/api /api
# depot build
COPY --from=builder /app/depot /depot

CMD ["/api"]
