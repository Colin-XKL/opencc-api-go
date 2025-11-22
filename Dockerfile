
FROM golang:1.25-alpine AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /opencc-api

FROM alpine AS deploy
WORKDIR /
COPY --from=build /opencc-api /opencc-api

# The original Dockerfile copied go modules to the final image.
# github.com/longbridgeapp/opencc might need dictionary files.
# Looking at the error log or library usage, it seems it embeds dictionaries or loads them.
# The library documentation says: `var Dir = flag.String("dir", defaultDir(), "dict dir")`
# But newer versions might use `embed`.
# The error `package embed is not in GOROOT` in the previous run suggested the code uses `embed`.
# If it uses `embed`, the dictionaries are inside the binary, so we might not need to copy /go/pkg/mod.
# However, to be safe and minimal changes, I will keep the structure but update the path if needed.
# Wait, `embed` puts files into the binary. So `COPY --from=build /go/pkg/mod/github.com ...` might be unnecessary if everything is embedded.
# But checking `github.com/longbridgeapp/opencc` source (I can't do that easily without cloning),
# let's assume `embed` is used for dictionaries.
# If I look at the `opencc-api.go` code, it imports the library.
# The previous Dockerfile copied `/go/pkg/mod/github.com`. This is very weird for a compiled Go binary unless it relies on files on disk at runtime.
# The `gwd0715/opencc` library likely needed external dictionary files.
# `longbridgeapp/opencc` uses `embed` (implied by the error requiring Go 1.16+).
# If it uses embed, we don't need to copy mod files.
# I'll try to remove the copy of pkg/mod and see. If it fails, I'll restore it.
# Actually, the user just asked to fix the runtime error and upgrade Go.
# The runtime error was in CI (GitHub Actions), not Docker (yet, but Docker would fail if it stays at 1.17 maybe? No, 1.17 supports embed).
# But the user wants 1.25.

COPY --from=build /go/pkg/mod/github.com /go/pkg/mod/github.com
EXPOSE 3000
LABEL author="Colin"
LABEL email="Colin_XKL@outlook.com"
ENTRYPOINT ["/opencc-api"]
