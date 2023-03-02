FROM golang:1.20.1-bullseye

WORKDIR /app

RUN apt-get update && apt-get install -y libzmq3-dev

COPY . ./

RUN go mod tidy
RUN go build

EXPOSE 50051

CMD [ "./moby-listener" ]