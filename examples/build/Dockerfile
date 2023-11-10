#
# Default values
#
ARG BINARY_NAME=publisher

#
# 1. Build Container
#
FROM golang:1.20 AS build
ARG VERSION
ARG CI_JOB_TOKEN
ARG CI_SERVER_HOST
ARG GITHUB_PAT

RUN bash -c 'echo -e "machine gitlab.com\nlogin gitlab-ci-token\npassword ${CI_JOB_TOKEN}\n" > ~/.netrc'
RUN git config --global url.https://gitlab-ci-token:${CI_JOB_TOKEN}@${CI_SERVER_HOST}.insteadOf https://${CI_SERVER_HOST}

RUN bash -c 'echo -e "machine github.com\nlogin devops-syntropy\npassword ${GITHUB_PAT}\n" >> ~/.netrc'
RUN git config --global url.https://devops-syntropy:${GITHUB_PAT}@github.com.insteadOf https://github.com

# Set up git credentials for private repositories
RUN mkdir -p ~/.netrc.d && \
    echo -e "machine gitlab.com\nlogin gitlab-ci-token\npassword ${CI_JOB_TOKEN}\n" > ~/.netrc.d/gitlab && \
    echo -e "machine github.com\nlogin devops-syntropy\npassword ${GITHUB_PAT}\n" > ~/.netrc.d/github && \
    chmod 600 ~/.netrc.d/*


# Copy necessary files for dependency resolution
WORKDIR /home/src
COPY go.mod go.sum ./

# Download dependencies
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# Copy the entire source code
COPY . .

# Build the binary
RUN make build

#
# 2. Runtime Container
#
FROM alpine as publisher

ENV TZ=Europe/Vilnius \
    PATH="/app:${PATH}"

# Install necessary packages and set timezone
RUN apk add --update --no-cache \
    tzdata \
    ca-certificates \
    bash \
    && \
    cp --remove-destination /usr/share/zoneinfo/${TZ} /etc/localtime && \
    echo "${TZ}" > /etc/timezone

# Copy the binary from the build container
WORKDIR /home/app
COPY --from=build /home/src/dist/${BINARY_NAME} .

CMD ["./${BINARY_NAME}", "--nats", "${NATS_HOSTS}", "--nats-jwt", "${NATS_JWT}", "--nats-nkey", "${NATS_NKEY}"]
