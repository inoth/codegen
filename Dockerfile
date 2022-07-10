FROM alpine@sha256:e7d88de73db3d3fd9b2d63aa7f447a10fd0220b7cbf39803c803f2af9ba256b3
LABEL maintainer = "inoth" version = "v1.0" description = "codegenfornet"
EXPOSE 8888
WORKDIR /

RUN apk add --no-cache bash
RUN apk --no-cache add ca-certificates
ENV GIN_MODE=release


COPY /release/. .

# RUN chmod +x ./setup.sh && chmod +x ./cmdb_notify
# ENTRYPOINT ["./setup.sh", "start"]
RUN chmod +x /codegenfornet
ENTRYPOINT ["/codegenfornet"]