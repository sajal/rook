# Copyright 2016 The Rook Authors. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM ubuntu:zesty
MAINTAINER Bassam Tabbara <bassam.tabbara@quantum.com>
WORKDIR /build
ENV TERM=xterm PATH=/usr/local/go/bin:/usr/local/osxcross/bin:$PATH

#
# install build tools, compilers and cross compilers for all supported platforms
#
RUN echo "deb [arch=amd64] http://archive.ubuntu.com/ubuntu zesty main universe restricted" > /etc/apt/sources.list &&\
    echo "deb [arch=amd64] http://archive.ubuntu.com/ubuntu zesty-updates main universe restricted" >> /etc/apt/sources.list &&\
    dpkg --add-architecture arm64 && \
    apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -yq --no-install-recommends \
        awscli \
        build-essential \
        ca-certificates \
        ccache \
        clang \
        cmake \
        crossbuild-essential-arm64 \
        curl \
        dwz \
        docker.io \
        git \
        jq \
        libxml2-dev \
        llvm-dev \
        mingw-w64 \
        patch \
        rsync \
        runit \
        sudo \
        uuid-dev \
        xz-utils \
        yasm \
        zip && \
    apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# install the OSX cross compilers. see https://github.com/tpoechtrager/osxcross
RUN OSXCROSS_COMMIT_HASH=86879571d33ff7252a55d777609742561dba9953 && \
    OSXCROSS_SDK=MacOSX10.11.sdk && \
    OSXCROSS_SDK_HASH=d75dddc8d820daf38e90a7c8e9771750a8a4371027a2d658a904cffacbf6b09b && \
    git clone https://github.com/quantum/osxcross && \
    git -C osxcross checkout -b castle-build ${OSXCROSS_COMMIT_HASH} && \
    curl -fsSL https://castle-cross-build.s3.amazonaws.com/osxcross/${OSXCROSS_SDK}.tar.xz -o osxcross/tarballs/${OSXCROSS_SDK}.tar.xz && \
    echo "${OSXCROSS_SDK_HASH}  osxcross/tarballs/${OSXCROSS_SDK}.tar.xz" | sha256sum -c - && \
    UNATTENDED=yes OSX_VERSION_MIN=10.6 osxcross/build.sh && \
    mv osxcross/target /usr/local/osxcross && \
    rm -rf osxcross && \
    rm -rf /usr/local/osxcross/SDK/${OSXCROSS_SDK}/usr/share

# install acbuild and docker2aci
RUN ACBUILD_VERSION=0.4.0 && \
    curl -fsSLO https://github.com/containers/build/releases/download/v${ACBUILD_VERSION}/acbuild-v${ACBUILD_VERSION}.tar.gz && \
    tar -C /usr/local/bin --strip-components=1 -xzf acbuild-v${ACBUILD_VERSION}.tar.gz && \
    rm -f acbuild-v${ACBUILD_VERSION}.tar.gz && \
    DOCKER2ACI_VERSION=0.14.0 && \
    curl -fsSLO https://github.com/appc/docker2aci/releases/download/v${DOCKER2ACI_VERSION}/docker2aci-v${DOCKER2ACI_VERSION}.tar.gz && \
    tar -C /usr/local/bin --strip-components=1 -xzf docker2aci-v${DOCKER2ACI_VERSION}.tar.gz && \
    rm -f docker2aci-v${DOCKER2ACI_VERSION}.tar.gz

# install go-lang and precompile the go standard library for all supported platforms and configurations
# there are three configurations that we compile for:
#    - default: CGO_ENABLED=1, defaults tags
#    - netgo: CGO_ENABLED=1, enable pure go DNS resolver
#    - nocgo: CGO_ENABLED=0, defaults tags, pure Go resolver is the default
#
RUN GO_VERSION=1.8 && \
    GO_HASH=53ab94104ee3923e228a2cb2116e5e462ad3ebaeea06ff04463479d7f12d27ca && \
    curl -fsSL https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz -o golang.tar.gz && \
    echo "${GO_HASH}  golang.tar.gz" | sha256sum -c - && \
    tar -C /usr/local -xzf golang.tar.gz && \
    rm golang.tar.gz && \
    \
    CGO_ENABLED=1 GOOS=darwin  GOARCH=amd64 CC=x86_64-apple-darwin15-clang   go install -a std && \
    CGO_ENABLED=1 GOOS=darwin  GOARCH=amd64 CC=x86_64-apple-darwin15-clang   go install -installsuffix netgo -tags netgo -a std && \
    CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 CC=x86_64-apple-darwin15-clang   go install -installsuffix nocgo -a std && \
    CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc        go install -a std && \
    CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc        go install -installsuffix netgo -tags netgo -a std && \
    CGO_ENABLED=0 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc        go install -installsuffix nocgo -a std && \
    CGO_ENABLED=1 GOOS=linux   GOARCH=arm64 CC=aarch64-linux-gnu-gcc         go install -a std && \
    CGO_ENABLED=1 GOOS=linux   GOARCH=arm64 CC=aarch64-linux-gnu-gcc         go install -buildmode pie -installsuffix pie -a std && \
    CGO_ENABLED=1 GOOS=linux   GOARCH=arm64 CC=aarch64-linux-gnu-gcc         go install -installsuffix netgo -tags netgo -a std && \
    CGO_ENABLED=0 GOOS=linux   GOARCH=arm64 CC=aarch64-linux-gnu-gcc         go install -installsuffix nocgo -a std && \
    CGO_ENABLED=1 GOOS=linux   GOARCH=amd64 CC=x86_64-linux-gnu-gcc          go install -a std && \
    CGO_ENABLED=1 GOOS=linux   GOARCH=amd64 CC=x86_64-linux-gnu-gcc          go install -buildmode pie -installsuffix pie -a std && \
    CGO_ENABLED=1 GOOS=linux   GOARCH=amd64 CC=x86_64-linux-gnu-gcc          go install -installsuffix netgo -tags netgo -a std && \
    CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 CC=x86_64-linux-gnu-gcc          go install -installsuffix nocgo -a std && \
    \
    CGO_ENABLED=1 GOOS=linux   GOARCH=amd64 CC=x86_64-linux-gnu-gcc   go install -installsuffix race -a std

COPY entrypoint.sh rsyncd.sh /build/
ENTRYPOINT [ "/build/entrypoint.sh" ]
