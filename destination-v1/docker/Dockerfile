FROM golang@sha256:fc801399d044a8e01f125eeb5aa3f160a0d12d6e03ba17a1d0b22ce50dfede81 AS builder

# Create the user and group files that will be used in the running container to
# run the process as an unprivileged user.
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

# Install the Certificate-Authority certificates for the app to be able to make
# calls to HTTPS endpoints.
# Git is required for fetching the dependencies.
RUN apk update \
    && apk add --no-cache git \
    && apk add ca-certificates

# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /src

# Fetch dependencies first; they are less susceptible to change on every build
# and will therefore be cached for speeding up the next build
COPY ./go.mod ./go.sum ./
RUN go mod download

# Import the code from the context.
COPY cmd/ cmd/
COPY data/ data/
COPY internals/ internals/

# Build the executable to `/app`. Mark the build as statically linked.
RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -o /app ./cmd/web/main.go

FROM scratch AS final

COPY --from=builder /user/group /user/passwd /etc/

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
# Import the compiled executable from the first stage.
COPY --from=builder /app /app

COPY --from=builder /src/data /data

EXPOSE 9001

USER nobody:nobody

ENTRYPOINT [ "/app" ]