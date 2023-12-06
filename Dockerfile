FROM golang:1.21 AS builder
WORKDIR /source/
COPY . /source/
RUN CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o bin/ ./...

FROM scratch
ENV CYBERMONDAY_ADDR ":80"
ENV CYBERMONDAY_DEFAULT_TEMPLATE_PATH "/usr/share/nginx/default/index.tplt.html"
ENV CYBERMONDAY_STATIC_DIR "/usr/share/nginx/html"
ENV CYBERMONDAY_BASEPATH "/usr/share/nginx/html"
ENV CYBERMONDAY_TITLE "Home"
VOLUME ["/usr/share/nginx/html/"]
ENTRYPOINT [ "/app/bin/cybermonday" ]
COPY resources/index.tplt.html /usr/share/nginx/default/
COPY resources/styles.css /usr/share/nginx/default/
COPY --from=builder /source/bin/* /app/bin/
