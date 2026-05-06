FROM golang:1.24-alpine AS builder

WORKDIR /src

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk add --no-cache git build-base

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

COPY api/go.mod api/go.sum ./api/
WORKDIR /src/api
RUN go mod download

COPY api/ ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/devops .

FROM alpine:3.20

WORKDIR /app

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk --no-cache add \
      ca-certificates \
      openssh-client \
      tzdata \
      curl \
      bash \
      git \
      kubectl

ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN mkdir -p /app/log /app/upload/xlsl /app/docs /app/common/templates

COPY --from=builder /out/devops /app/devops
COPY --from=builder /src/api/common/templates /app/common/templates
COPY --from=builder /src/api/upload/xlsl/host.xlsx /app/upload/xlsl/host.xlsx

RUN chmod +x /app/devops

EXPOSE 8000

CMD ["/app/devops"]
