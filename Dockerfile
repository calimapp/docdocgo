FROM scratch

COPY docdocgo /app/docdocgo

ENTRYPOINT ["/app/docdocgo"]