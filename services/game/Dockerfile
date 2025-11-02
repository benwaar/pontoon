FROM golang:1.22 AS build
WORKDIR /app
COPY . .
RUN go build -o game-service .

FROM gcr.io/distroless/base-debian12
COPY --from=build /app/game-service /game-service
EXPOSE 9000
ENTRYPOINT ["/game-service"]
