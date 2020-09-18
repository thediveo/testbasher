ARG GOVERSION=1.15

FROM golang:${GOVERSION}
WORKDIR /src
# Cache dependencies to some extend, to speed up things...
COPY go.mod go.sum ./
RUN go mod graph | awk '{if ($1 !~ "@") print $2}' | xargs go get
# Copy in the full module; unfortunately, tests cannot be prebuild, so that's
# all we can do here.
COPY . .
# Tests must be started as root, so we can run the whole test suite once as root
# and another time as non-root.
CMD ["go", "test", "-v", "."]
