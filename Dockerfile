FROM golang:1.12-alpine AS builder
RUN apk add --no-cache git
ENV GO111MODULE on
COPY go.* /go/github.com/gkawamoto/cybermonday/
COPY cmd /go/github.com/gkawamoto/cybermonday/cmd
RUN mkdir -p /result/ && cd /go/github.com/gkawamoto/cybermonday/ &&\
 go build -o /result/entrypoint cmd/entrypoint/main.go &&\
 go build -o /result/cybermonday cmd/cybermonday/main.go

FROM nginx:stable-alpine
RUN rm -v /usr/share/nginx/html/*
ENV CYBERMONDAY_BASEPATH /usr/share/nginx/html
VOLUME ["/usr/share/nginx/html/"]
ENTRYPOINT [ "/usr/bin/entrypoint" ]
COPY nginx/default.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /result/* /usr/bin/
