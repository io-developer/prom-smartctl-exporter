FROM debian:stretch-slim

RUN apt-get update \
    && apt-get install smartmontools -y --no-install-recommends \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

ADD bin/prom-smartctl-exporter /prom-smartctl-exporter

EXPOSE 9167
CMD ["/prom-smartctl-exporter"]
