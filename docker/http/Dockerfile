FROM alpine:3.5

RUN mkdir -p /bin/needle-url/

ADD server /bin/needle-url/server
ADD config.json /bin/needle-url/config.json

WORKDIR "/bin/needle-url"
CMD ["./server"]