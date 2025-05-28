FROM debian:12.11-slim



# Create a user that will be used when inisde the container rather than root
RUN useradd dev --create-home && \
    echo "dev:pass" | chpasswd && \
    adduser dev sudo

# bash-completion, bsdmainutils (column), gcc, jq, make and vim are useful development tools
# curl is needed for installation of Go and other tools
# git is needed to run commands such as go mod tidy during development
# postgresql-client install psql for connecting to postgres databases
# sudo is useful to have inside development container, to add missing tools
# wget is used to download tools
RUN apt-get update && apt-get -y install \
    bash-completion \
    bsdmainutils \
    curl \
    gcc \
    git \
    jq \
    make \
    postgresql-client \
    sudo \
    vim \
    wget

# Install Go
ARG GO_VERSION=1.24.3.linux-amd64
# Modify PATH to make go immediately available after installation
ENV PATH=$PATH:/usr/local/go/bin
RUN curl -sSfL "https://go.dev/dl/go${GO_VERSION}.tar.gz" | tar -C /usr/local -xzf -  && \
    go env -w GOBIN=/usr/local/go/bin GOROOT=/usr/local/go && \
    echo "export GOROOT=/usr/local/go" > /etc/profile.d/go.sh  && \
    echo "export PATH=\$PATH:/usr/local/go/bin" >> /etc/profile.d/go.sh

# Install Go linter
ARG GOLANGCI_VERSION=v2.1.6
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOROOT)/bin ${GOLANGCI_VERSION}

# Instal goimport-reviser
ARG GOIMPORTS_REVISER_VERSION=v3.9.1
RUN go install -v github.com/incu6us/goimports-reviser/${GOIMPORTS_REVISER_VERSION%%.*}@${GOIMPORTS_REVISER_VERSION}

# Install mockery
ARG MOCKERY_VERSION=v3.3.0/mockery_3.3.0_Linux_x86_64
RUN wget -qO- https://github.com/vektra/mockery/releases/download/${MOCKERY_VERSION}.tar.gz | tar -xzv -C /usr/local/bin mockery

# Install staticcheck
ARG STATICCHECK_VERSION=2025.1.1/staticcheck_freebsd_amd64
RUN wget -qO- https://github.com/dominikh/go-tools/releases/download/${STATICCHECK_VERSION}.tar.gz | tar -xzv -C /usr/local/bin --strip-components 1 staticcheck/staticcheck

# Install gopls
ARG GOPLS_VERSION=v0.18.1
RUN go install golang.org/x/tools/gopls@${GOPLS_VERSION}

# Install pnpm (useful tools are written in node)
# Pnpm does not support installing once for multiple users, so install for the dev user and switch back to the root user after
USER dev
WORKDIR /home/dev
ARG PNPM_VERSION=10.11.0
ARG NODE_VERSION=22
RUN wget -qO- https://get.pnpm.io/install.sh | PNPM_VERSION=${PNPM_VERSION} ENV="$HOME/.bashrc" SHELL="$(which bash)" bash -
RUN export PNPM_HOME=/home/dev/.local/share/pnpm && \
    export PATH=$PATH:$PNPM_HOME && \
    pnpm env use --global ${NODE_VERSION}
USER root
WORKDIR /

# Install openapi-schema-validator
# Pnpm does not support installing for multiple users, so install for the dev user and switch back to the root user after
USER dev
WORKDIR /home/dev
ARG OPENAPI_SCHEMA_VALIDATOR_VERSION=2.4.1
RUN export PNPM_HOME=/home/dev/.local/share/pnpm && \
    export PATH=$PATH:$PNPM_HOME && \
    pnpm install -g @seriousme/openapi-schema-validator@${OPENAPI_SCHEMA_VALIDATOR_VERSION}
USER root
WORKDIR /

CMD ["sh"]
ENTRYPOINT ["bash", "-lec"]
