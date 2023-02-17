FROM golang:buster AS builder
# Debian 10 O.S., Linux Kernel 4.19 (LTS).
#TODO: Update to 'golang:bullseye' when it is released.
# Debian 11 O.S., Linux Kernel 5.10 (LTS).

# Show version of important things.
RUN echo "$PATH"
RUN uname --all
RUN go version
RUN go env
RUN openssl version

# Update the packages and C.A. certificates.
RUN DEBIAN_FRONTEND=noninteractive \
  apt-get -y update && \
  apt-get -y upgrade && \
  update-ca-certificates

# Where are we ?
RUN echo $PWD
RUN ls -l

# Copy the source code.
ARG SOURCE_FOLDER=/go/src/github.com/vault-thirteen/gfe
RUN echo "Source Folder = '$SOURCE_FOLDER'"
RUN mkdir -p "$SOURCE_FOLDER"
COPY . "$SOURCE_FOLDER"

# Build the application.
RUN mkdir -p "/app"
WORKDIR "$SOURCE_FOLDER"
RUN ls -l
RUN go mod vendor
RUN cd "cmd/server" && go build -o server && mv "./server" "/app/server"

# While the Vault is not available,
# Prepare a temporary public key for JWT tokens.
WORKDIR "$SOURCE_FOLDER/test/web-token"
RUN go install github.com/golang-jwt/jwt/cmd/jwt@v3.2.1
RUN ls -l && rm -rf "./key" && mkdir "key" && ls -l
RUN ./create-rsa-key-and-token.sh quiet=yes
RUN mkdir -p "/app/data/key"
RUN cp "./key/jwt_rsa_RS256.key" "/app/data/key/jwt_rsa_RS256.key"
RUN cp "./key/jwt_rsa_RS256.key.pub" "/app/data/key/jwt_rsa_RS256.key.pub"
RUN cp "./key/jwt_token.txt" "/app/data/key/jwt_token.txt"

FROM debian:buster AS runner
# Debian 10 O.S., Linux Kernel 4.19 (LTS).
#TODO: Update to 'debian:bullseye' when it is released.
# Debian 11 O.S., Linux Kernel 5.10 (LTS).

# Show version of important things.
RUN echo "$PATH"
RUN uname --all

# Update the packages and C.A. certificates.
RUN DEBIAN_FRONTEND=noninteractive \
  apt-get -y update && \
  apt-get -y upgrade && \
  apt-get install -y ca-certificates
RUN apt-get install -y nano mc htop iotop

# Where are we ?
RUN echo $PWD
RUN ls -l

# Copy files from builder.
RUN mkdir -p "/app/data/key"
COPY --from=builder "/app/server" "/app/server"
COPY --from=builder "/app/data/key/jwt_rsa_RS256.key" "/app/data/key/jwt_rsa_RS256.key"
COPY --from=builder "/app/data/key/jwt_rsa_RS256.key.pub" "/app/data/key/jwt_rsa_RS256.key.pub"
COPY --from=builder "/app/data/key/jwt_token.txt" "/app/data/key/jwt_token.txt"

# Set the environment variables.
ENV GFE_IS_DEBUG_ENABLED="true"
ENV GFE_JWT_KEY_SOURCE_TYPE="2"
ENV GFE_JWT_KEY_DSN="file:///app/data/key/jwt_rsa_RS256.key.pub"
ENV GFE_JWT_KEY_VALUE=""
ENV GFE_KAFKA_CONSUMER_GROUP_ID="gfe_test_consumer_group"
ENV GFE_KAFKA_BROKER_ADDRESS_LIST="127.0.0.1:9093"
ENV GFE_KAFKA_TOPIC_LIST="topic_a"
ENV GFE_POSTGRE_USER="test"
ENV GFE_POSTGRE_PASSWORD="test"
ENV GFE_POSTGRE_DATABASE="test"
ENV GFE_METRICS_HTTP_SERVER_PORT="9999"
ENV GFE_BUSINESS_HTTP_SERVER_PORT="8080"

# Start the application.
WORKDIR "/app"
ENTRYPOINT ["/app/server"]

# Firewall.
EXPOSE 8080/TCP
EXPOSE 9999/TCP
