FROM golang:1.10

ARG workdir=github.com/edrans/cloudflarebeat


ENV WORKDIR=${workdir}
ENV VERSION=latest

COPY build/build_os_binaries.sh /usr/bin

RUN curl https://glide.sh/get | sh
RUN chmod 755 /usr/bin/build_os_binaries.sh
RUN mkdir -p /go/src/${workdir}

WORKDIR /go/src/${workdir}
COPY . .

CMD ["/usr/bin/build_os_binaries.sh"]
