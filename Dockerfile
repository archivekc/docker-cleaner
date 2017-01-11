FROM alpine:latest

COPY docker-cleaner /bin/docker-cleaner

CMD /bin/docker-cleaner
