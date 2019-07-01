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
ENV CYBERMONDAY_BOOTSTRAP_REF //stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css
ENV CYBERMONDAY_TITLE Home
ENV PATH $PATH:/application/bin/
VOLUME ["/usr/share/nginx/html/"]
ENTRYPOINT [ "/application/bin/entrypoint" ]
COPY nginx/default.conf /etc/nginx/conf.d/default.conf
COPY resources/default.template.html /application/default.template.html
COPY --from=builder /result/* /application/bin/
