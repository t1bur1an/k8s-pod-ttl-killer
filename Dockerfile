FROM scratch

ARG BIN_FILE=k8s-pod-ttl-killer

ADD ${BIN_FILE} /

ENTRYPOINT ["/${BIN_FILE}"]
