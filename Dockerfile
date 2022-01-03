
FROM golang:1.17-alpine AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /opencc-api

FROM alpine AS deploy
WORKDIR /
COPY --from=build /opencc-api /opencc-api

COPY --from=build /go/pkg/mod/github.com /go/pkg/mod/github.com
EXPOSE 8000
LABEL author="Colin"
LABEL email="Colin_XKL@outlook.com"
ENTRYPOINT ["/opencc-api"]