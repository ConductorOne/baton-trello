FROM gcr.io/distroless/static-debian11:nonroot
ENTRYPOINT ["/baton-trello"]
COPY baton-trello /