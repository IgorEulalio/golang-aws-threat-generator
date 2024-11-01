FROM alpine:3.18

COPY main /
RUN chmod +x /main

USER 1004

ENTRYPOINT ["/main"]
