FROM golang:1.26 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/dota .

FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=build /out/dota .
COPY templates ./templates
EXPOSE 8080
ENTRYPOINT ["/app/dota"]
