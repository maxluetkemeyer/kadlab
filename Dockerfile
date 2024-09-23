FROM golang:alpine

RUN go install github.com/vadimi/grpc-client-cli/cmd/grpc-client-cli@latest

WORKDIR /kademlia

COPY . /kademlia

COPY /proto /etc/proto

RUN go build -o kademlia-node

# used to debug grpc calls
RUN echo "/go/bin/grpc-client-cli --proto /etc/proto/kademlia.proto \$@" >> /bin/grpc

RUN chmod +x /bin/*

EXPOSE 50051

ENTRYPOINT [ "/kademlia/kademlia-node" ]