
FROM golang:1.25-alpine AS build
WORKDIR /app
COPY . .
RUN go mod download

ARG VERSION
ARG BUILD_TIME
ARG GIT_COMMIT

RUN CGO_ENABLED=0 go build -ldflags="-s -w -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}" -o /opencc-api

FROM alpine AS deploy
WORKDIR /
COPY --from=build /opencc-api /opencc-api

EXPOSE 3000
LABEL author="Colin"
ENTRYPOINT ["/opencc-api"]
