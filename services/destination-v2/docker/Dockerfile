FROM golang:alpine AS builder

RUN mkdir /root/builder

# Create the user and group files that will be used in the running container to
# run the process as an unprivileged user.
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /root/builder

# Fetch dependencies first; they are less susceptible to change on every build
# and will therefore be cached for speeding up the next build
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . /root/builder
RUN apk update \
 && apk add git ca-certificates \
 && cd /root/builder \
 && go build -o ../service .

FROM alpine:latest AS final

COPY --from=builder /user/group /user/passwd /etc/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /root/service /app
EXPOSE 9001
USER nobody:nobody
ENTRYPOINT [ "/app" ]
