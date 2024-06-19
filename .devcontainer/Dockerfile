ARG USERNAME=user
ARG USER_UID=1000
ARG USER_GID=$USER_UID 
ARG DEBIAN_VERSION=bullseye
ARG GO_VERSION=1.22
ARG GOMODIFYTAGS_VERSION=v1.16.0
ARG GOPLAY_VERSION=v1.0.0
ARG GOTESTS_VERSION=v1.6.0
ARG DLV_VERSION=v1.22.0
ARG MOCKERY_VERSION=v2.40.1
ARG GOMOCK_VERSION=v1.6.0
ARG MOCKGEN_VERSION=v1.6.0
ARG GOPLS_VERSION=v0.14.2
ARG GOLANGCILINT_VERSION=v1.56.2
ARG IMPL_VERSION=v1.2.0
ARG GOPKGS_VERSION=v2.1.2
ARG KUBECTL_VERSION=v1.29.1
ARG STERN_VERSION=v1.28.0
ARG KUBECTX_VERSION=v0.9.5
ARG KUBENS_VERSION=v0.9.5
ARG HELM_VERSION=v3.14.0


FROM golang:${GO_VERSION}-${DEBIAN_VERSION} AS go
FROM iznth/binbase:gomodifytags-${GOMODIFYTAGS_VERSION} AS gomodifytags
FROM iznth/binbase:goplay-${GOPLAY_VERSION} AS goplay
FROM iznth/binbase:gotests-${GOTESTS_VERSION} AS gotests
FROM iznth/binbase:dlv-${DLV_VERSION} AS dlv
FROM iznth/binbase:mockery-${MOCKERY_VERSION} AS mockery
FROM iznth/binbase:gomock-${GOMOCK_VERSION} AS gomock
FROM iznth/binbase:mockgen-${MOCKGEN_VERSION} AS mockgen
FROM iznth/binbase:gopls-${GOPLS_VERSION} AS gopls
FROM iznth/binbase:golangci-lint-${GOLANGCILINT_VERSION} AS golangci-lint
FROM iznth/binbase:impl-${IMPL_VERSION} AS impl
FROM iznth/binbase:gopkgs-${GOPKGS_VERSION} AS gopkgs
FROM iznth/binbase:kubectl-${KUBECTL_VERSION} AS kubectl
FROM iznth/binbase:stern-${STERN_VERSION} AS stern
FROM iznth/binbase:kubectx-${KUBECTX_VERSION} AS kubectx
FROM iznth/binbase:kubens-${KUBENS_VERSION} AS kubens
FROM iznth/binbase:helm-${HELM_VERSION} AS helm


FROM iznth/baseddevcontainer:debian as base


ARG CREATED
ARG COMMIT
ARG VERSION=local
LABEL \
    org.opencontainers.image.authors="brad@izn.ai" \
    org.opencontainers.image.created=$CREATED \
    org.opencontainers.image.version=$VERSION \
    org.opencontainers.image.revision=$COMMIT \
    org.opencontainers.image.url="https://github.com/iznth/godevcontainer" \
    org.opencontainers.image.documentation="https://github.com/iznth/godevcontainer" \
    org.opencontainers.image.source="https://github.com/iznth/godevcontainer" \
    org.opencontainers.image.title="Go Dev container Debian" \
    org.opencontainers.image.description="Go development container for Visual Studio Code Remote Containers development"
COPY --from=go /usr/local/go /usr/local/go
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH \
    CGO_ENABLED=0 \
    GO111MODULE=on
WORKDIR $GOPATH
# Install Debian packages
RUN apt-get update && \
    apt-get install -y --no-install-recommends g++ wget && \
    apt-get autoremove -y && \
    apt-get clean -y && \
    rm -r /var/cache/* /var/lib/apt/lists/* 

COPY --from=gomodifytags /bin /go/bin/gomodifytags
COPY --from=goplay  /bin /go/bin/goplay
COPY --from=gotests /bin /go/bin/gotests
COPY --from=dlv /bin /go/bin/dlv
COPY --from=mockery /bin /go/bin/mockery
COPY --from=gomock /bin /go/bin/gomock
COPY --from=mockgen /bin /go/bin/mockgen
COPY --from=gopls /bin /go/bin/gopls
COPY --from=golangci-lint /bin /go/bin/golangci-lint
COPY --from=impl /bin /go/bin/impl
COPY --from=gopkgs /bin /go/bin/gopkgs

# Extra binary tools
COPY --from=kubectl /bin /usr/local/bin/kubectl
COPY --from=stern /bin /usr/local/bin/stern
COPY --from=kubectx /bin /usr/local/bin/kubectx
COPY --from=kubens /bin /usr/local/bin/kubens
COPY --from=helm /bin /usr/local/bin/helm

FROM base as dev

# # Create the user
# RUN groupadd --gid $USER_GID $USERNAME
# RUN useradd --uid $USER_UID --gid $USER_GID -m $USERNAME
# #
# # [Optional] Add sudo support. Omit if you don't need to install software after connecting.
# RUN apt-get update \
#     && apt-get install -y sudo
# RUN echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME
# RUN chmod 0440 /etc/sudoers.d/$USERNAME

# # [Optional] Set the default user. Omit if you want to keep the default as root.
# USER $USERNAME
