FROM alpine

ADD aws_status_linux_amd64 /
ADD config.yml /

RUN apk --update upgrade && \
    apk add curl ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/*

EXPOSE 8080

CMD /aws_status_linux_amd64
