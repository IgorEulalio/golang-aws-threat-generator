# FROM public.ecr.aws/o5z5d3z9/pgcrooks/workload-agent:gcr-beta as workload-agent

FROM alpine:3.18

COPY main /
RUN chmod +x /main

COPY --from=workload-agent /opt/draios /opt/draios

ENV SYSDIG_WORKLOAD_ID=release-book-gcr-test

ENTRYPOINT ["/opt/draios/bin/instrument", "/main"]
# ENTRYPOINT ["/main"]
