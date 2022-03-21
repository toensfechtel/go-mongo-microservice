FROM public.ecr.aws/amazonlinux/amazonlinux:latest as build

ENV GOPATH /go
#Installs go 
RUN yum install -y golang
#A GOPROXY controls the source of your Go module downloads.
RUN go env -w GOPROXY=direct
ENV PATH="${GOPATH}/src:${GOPATH}/src/bin:${PATH}"
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