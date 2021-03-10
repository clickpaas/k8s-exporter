FROM alpine  as builder
ARG VERSION="v20210310-1"
ENV GOPROXY https://goproxy.cn
RUN apk update && apk add git go && mkdir /root/go/src -pv && \
    cd /root/go/src/ && \
    git clone https://github.com/clickpaas/k8s-exporter.git && \
    cd k8s-exporter && go mod vendor && go build .

RUN cd /root/go/src && \
    git clone https://github.com/go-delve/delve && \
    cd delve/cmd/dlv && go build .


FROM alpine
WORKDIR /root
COPY --from=builder /root/go/src/k8s-exporter/clickpaas-exporter .
COPY --from=builder /root/go/src/delve/cmd/dlv /usr/bin/dlv

RUN apk update && apk add  bash curl tzdata && \
    cp -r -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    chmod +x clickpaas-exporter

ENTRYPOINT ["./clickpaas-exporter"]
