FROM scratch
ADD ./ssl/ca-certificates.crt /etc/ssl/certs/
ADD ./build/distate-task /
CMD ["/distate-task"]
