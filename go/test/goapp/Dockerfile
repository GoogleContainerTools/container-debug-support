ARG GOVERSION
FROM golang:$GOVERSION as builder
ARG BUILDPLATFORM=linux/arm64
ARG TARGETOS=linux
ARG TARGETARCH=arm64

COPY main.go .
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -gcflags="all=-N -l" -o /app main.go

FROM gcr.io/distroless/base
CMD ["./app"]
ENV GOTRACEBACK=single
COPY --from=builder /app .
