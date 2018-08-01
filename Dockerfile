FROM alpine:latest
# Now just add the binary
COPY app /
COPY swagger.json /
ENTRYPOINT ["/app"]
EXPOSE 8005
