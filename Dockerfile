FROM busybox:latest

ARG BIN_FILE=k8s-pod-ttl-killer

WORKDIR /app

COPY ${BIN_FILE} ./svc

ENTRYPOINT ["./svc"]
