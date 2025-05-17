FROM ubuntu:22.04

RUN apt-get update && \
    apt-get install -y spatialite-bin && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /data

ENTRYPOINT ["spatialite"]

CMD []