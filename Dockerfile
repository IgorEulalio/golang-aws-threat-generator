FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y \
    bash \
    curl \
    unzip \
    && rm -rf /var/lib/apt/lists/*

# Install AWS CLI v2
RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-$(uname -m).zip" -o "/tmp/awscliv2.zip" \
    && unzip /tmp/awscliv2.zip -d /tmp \
    && /tmp/aws/install \
    && rm -rf /tmp/aws*

COPY main /
COPY trigger-rules.sh /

RUN chmod +x /trigger-rules.sh
RUN chmod +x /main

USER 1004

ENTRYPOINT ["/main"]
