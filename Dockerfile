FROM golang:1.22.0

WORKDIR /server-files

COPY . .

RUN apt update && apt install -y zip 

RUN go build -o server.out .

EXPOSE 4000

CMD ["./server.out"]
