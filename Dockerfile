FROM golang:1.25.4-alpine3.22 AS build

WORKDIR /app

COPY go.* .
COPY cmd/ ./cmd/
COPY internal/ ./internal/

RUN go build -o /bin/controller ./cmd/controller/ 


FROM gcr.io/distroless/static

COPY --from=build /bin/controller /controller

ENTRYPOINT ["/controller"]
