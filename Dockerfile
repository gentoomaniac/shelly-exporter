FROM scratch

COPY shelly-exporter /

ENTRYPOINT ["/shelly-exporter"]
