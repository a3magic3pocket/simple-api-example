FROM golang:1.16

WORKDIR /root/api

COPY . /root/api

RUN cp /etc/apt/sources.list sources.list_backup
RUN sed -i -e 's/archive.ubuntu.com/mirror.kakao.com/g' /etc/apt/sources.list
RUN sed -i -e 's/security.ubuntu.com/mirror.kakao.com/g' /etc/apt/sources.list

RUN go mod tidy
RUN go build -o api
RUN mv .env.example .env

ENTRYPOINT ./api