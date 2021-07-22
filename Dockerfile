FROM ubuntu:20.10

ENV EXPORTER_HOST="0.0.0.0"
ENV EXPORTER_PORT="9167"
ENV EXPORTER_SHELL_TEMPLATE="%s"

RUN apt-get update \
    && apt-get install smartmontools -y --no-install-recommends \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

ADD bin/prom-smartctl-exporter /prom-smartctl-exporter

EXPOSE ${EXPORTER_PORT}

CMD ["/prom-smartctl-exporter", "-listen", ${EXPORTER_HOST}:${EXPORTER_PORT}, "-shell", ${EXPORTER_SHELL_TEMPLATE}]
