FROM ubuntu:latest
RUN mkdir -p /chto-tam-po-peresdacham
WORKDIR /chto-tam-po-peresdacham
COPY chto-tam-po-peresdacham .
ENTRYPOINT [ "./chto-tam-po-peresdacham" ]
