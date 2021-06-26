FROM ubuntu:bionic 

WORKDIR /go/bin/
COPY ./hakstore-server .
COPY ./config.yml .

RUN apt update
RUN apt install -y ca-certificates

CMD ["./hakstore-server"]
