FROM scratch
ADD rootcerts.crt /etc/ssl/certs/ca-certificates.crt
ADD swarm /swarm
ENTRYPOINT ["/swarm"]
EXPOSE 2375
CMD ["-h"]
