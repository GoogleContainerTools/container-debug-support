ARG GOVERSION
FROM --platform=$BUILDPLATFORM golang:$GOVERSION as builder
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

COPY main.go .
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -gcflags="all=-N -l" -o /app main.go

FROM --platform=$BUILDPLATFORM gcr.io/distroless/base
CMD ["./app"]
ENV GOTRACEBACK=single
COPY --from=builder /app .
