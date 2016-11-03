FROM busybox
ENV LISTEN_PORT 8888
EXPOSE 8888
ADD authz-webhook /
CMD ["/authz-webhook"]
