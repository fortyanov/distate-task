FROM golang:1.15-alpine

#WORKDIR /app/distate-task
WORKDIR /go/src/distate-task

COPY . .

#RUN mkdir /app/conf
#RUN cp /app/distate-task/src/tern.conf /app/conf/tern.conf
#COPY ./tern.conf /app/conf/tern.conf

RUN go get -d -v ./...
RUN go install -v ./...

#RUN go get -u github.com/jackc/tern
#RUN tern migrate --migrations /home/forty/projects/distate-task/migrations
RUN go build -o /go/out/distate-task .
RUN rm -rf /go/src/distate-task

EXPOSE 8080

CMD ["/go/out/distate-task"]
