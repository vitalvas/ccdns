FROM ubuntu:latest
RUN apt update -qy && apt install -qy ca-certificates ssl-cert
COPY ccdns /bin/ccdns
CMD ["/bin/ccdns"]
