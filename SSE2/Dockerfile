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
ARG SOURCE_FOLDER=/go/src/github.com/vault-thirteen/SSE2
RUN echo "Source Folder = '$SOURCE_FOLDER'"
RUN mkdir -p "$SOURCE_FOLDER"
#COPY [--chown=<user>:<group>] <src>... <dest>
COPY . "$SOURCE_FOLDER"

# Build the application.
RUN mkdir -p "/app"
WORKDIR "$SOURCE_FOLDER"
RUN ls -l
RUN go mod vendor
RUN cd "cmd/server" && go build -o server && mv "./server" "/app/server"

#FROM debian:buster AS runner
# Debian 10 O.S., Linux Kernel 4.19 (LTS).
#TODO: Update to 'debian:bullseye' when it is released.
# Debian 11 O.S., Linux Kernel 5.10 (LTS).
FROM localhost:5000/libre_office AS runner

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
ARG SOURCE_FOLDER=/go/src/github.com/vault-thirteen/SSE2
RUN echo "Source Folder = '$SOURCE_FOLDER'"
RUN mkdir -p "/app/configs/server"
RUN mkdir -p "/app/data"
#COPY [--chown=<user>:<group>] <src>... <dest>
COPY --from=builder "/app/server" "/app/server"
COPY --from=builder "$SOURCE_FOLDER/configs/server/file_size_limits.xml" "/app/configs/server/file_size_limits.xml"

# Set the environment variables.
ENV SSE2_IS_DEBUG_ENABLED="true"
ENV SSE2_WORKERS_COUNT="2"
ENV SSE2_PATH_TO_CONVERTER_EXECUTABLE="soffice"
ENV SSE2_LARGE_PNG_IMAGE_MAXIMUM_SIDE_DIMENSION="640"
ENV SSE2_SMALL_PNG_IMAGE_MAXIMUM_SIDE_DIMENSION="240"
ENV SSE2_FILE_SIZE_LIMIT_SETTINGS_FILE="configs/server/file_size_limits.xml"
ENV SSE2_USE_LIBRE_OFFICE_MULTIPLE_USER_INSTALLATIONS=true
ENV SSE2_S3_SERVER_ADDRESS="http://127.0.0.1:9002"
ENV SSE2_S3_ACCESS_KEY="serious"
ENV SSE2_S3_SECRET="shit"
ENV SSE2_S3_TOKEN="token"
ENV SSE2_S3_REGION="region"
#ENV SSE2_S3_DISABLE_SSL=""
#ENV SSE2_S3_FORCE_PATH_STYLE=""
ENV SSE2_S3_IS_MINIO="true"
ENV SSE2_S3_LOCAL_FILES_FOLDER="/app/data"
ENV SSE2_INPUT_KAFKA_CONSUMER_GROUP_ID="sse2_test_consumer_group"
ENV SSE2_INPUT_KAFKA_BROKER_ADDRESS_LIST="127.0.0.1:9093"
ENV SSE2_INPUT_KAFKA_TOPIC_LIST="sse2_tasks"
ENV SSE2_OUTPUT_KAFKA_BROKER_ADDRESS_LIST="127.0.0.1:9093"
ENV SSE2_OUTPUT_KAFKA_TOPIC_LIST="sse2_task_results"
#ENV SSE2_HTTP_SERVER_HOST=""
ENV SSE2_HTTP_SERVER_PORT="8080"

# Start the application.
WORKDIR "/app"
ENTRYPOINT ["/app/server"]

# Firewall.
EXPOSE 8080/TCP
