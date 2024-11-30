# Build container
FROM golang:bookworm AS build

# Setup environment
RUN mkdir -p /data
WORKDIR /data

# Install native dependencies
RUN apt update
RUN apt install -y protobuf-compiler

# Build the release
COPY . .
RUN make depend
RUN make build/cli

# Extract the release
RUN mkdir -p /out
RUN cp out/bofied-backend /out/bofied-backend

# Release container
FROM debian:bookworm

# Add certificates
RUN apt update
RUN apt install -y ca-certificates

# Add the release
COPY --from=build /out/bofied-backend /usr/local/bin/bofied-backend

CMD /usr/local/bin/bofied-backend
