FROM debian:bookworm-slim

RUN groupadd -r sandbox && useradd -r -g sandbox sandbox

COPY lango /usr/local/bin/lango

USER sandbox

ENTRYPOINT ["/usr/local/bin/lango", "--sandbox-worker"]
