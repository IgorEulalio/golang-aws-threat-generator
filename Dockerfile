FROM alpine:3.18

COPY main /
RUN chmod +x /main

USER 1000

ENTRYPOINT ["/main"]
