
FROM golang:1.25-alpine AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /opencc-api

FROM alpine AS deploy
WORKDIR /
COPY --from=build /opencc-api /opencc-api

EXPOSE 3000
LABEL author="Colin"
ENTRYPOINT ["/opencc-api"]
