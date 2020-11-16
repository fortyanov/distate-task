#CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/distate-task .
#docker build -t distate-task -f Dockerfile .

FROM scratch
#FROM alpine
ADD ssl /etc/ssl/certs/
ADD build/distate-task /
CMD ["/distate-task"]