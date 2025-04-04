FROM gcr.io/distroless/static-debian11:nonroot
ENTRYPOINT ["/baton-fluid-topics"]
COPY baton-fluid-topics /