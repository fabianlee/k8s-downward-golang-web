#
# builder image
#
FROM golang:1.16.3-buster as builder
RUN mkdir /build
ADD src/*.go /build/
WORKDIR /build

# accept override of value from --build-args
ARG MY_VERSION=0.1.1
ENV MY_VERSION=$MY_VERSION

# accept override of value from --build-args
ARG MY_BUILDTIME=now
ENV MY_BUILDTIME=$MY_BUILDTIME

# create module, fetch dependencies, then build
RUN go mod init fabianlee.org/k8s-downward-golang-web \
   && go get -d -u ./... \
   && CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.Version=${MY_VERSION} -X main.BuildTime=${MY_BUILDTIME}" -a -o fabianlee.org/main .


#
# generate small final image for end users
#
#FROM alpine:3.13.5
# could have used either alpine or busybox
# busybox-glibc (versus musl) would have better compatability with Debian, but that is not an issue here
FROM busybox:1.32.1-glibc

# copy golang binary into container
WORKDIR /root
COPY --from=builder /build/fabianlee.org/main .

# executable
ENTRYPOINT [ "./main" ]
