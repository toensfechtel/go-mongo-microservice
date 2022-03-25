FROM public.ecr.aws/amazonlinux/amazonlinux:latest as build


ENV GOPATH /go

ENV GOLANG_VERSION 1.17.8
ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz
ENV GOLANG_DOWNLOAD_SHA256 980e65a863377e69fd9b67df9d8395fd8e93858e7a24c9f55803421e453f4f99

RUN yum update -y 

RUN yum install -y yum-utils make gcc gcc-c++ ghostscript git git-core tar

# install go
RUN curl -fsSL "${GOLANG_DOWNLOAD_URL}" -o golang.tar.gz \
  && echo "${GOLANG_DOWNLOAD_SHA256}  golang.tar.gz" | sha256sum -c - \
  && tar -C /usr/local -xzf golang.tar.gz \
  && rm golang.tar.gz

ENV PATH $GOPATH/bin:/usr/local/go/bin:$GOPATH/src:$GOPATH/src/bin:$PATH

RUN go env -w GOPROXY=direct

WORKDIR $GOPATH/src
#go-delve is a tool to help us debug inside the container, if you 
#have chosen not to use the amazon linux image make sure to install this
RUN go get github.com/go-delve/delve/cmd/dlv
ENV PATH="${GOPATH}/src:${GOPATH}/src/bin:${GOPATH}/src/github.com/go-delve/delve:${GOPATH}/bin:${PATH}"
#Copies project over
COPY . .
#Builds Gongo
RUN go build 
ENV GIN_MODE release
ENV HOST 0.0.0.0
ENV PORT 8080
EXPOSE 8080 2345

CMD ["gongo"]