FROM golang:1.22.5-alpine as build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN GOOS=linux go build -o /feedback-api


FROM gcr.io/distroless/static-debian12

WORKDIR /

COPY --from=build-stage /feedback-api /feedback-api

ENV API_LISTEN_PORT 8080
EXPOSE 8080
EXPOSE 2222

USER nonroot:nonroot

ENTRYPOINT ["/feedback-api"]
