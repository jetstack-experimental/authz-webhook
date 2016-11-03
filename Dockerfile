FROM busybox
ADD bitesize-authz-webhook /
CMD ["/bitesize-authz-webhook"]
