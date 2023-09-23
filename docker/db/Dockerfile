FROM postgres:15

# Postgresのロケールを日本語に設定
RUN localedef -i ja_JP -c -f UTF-8 -A /usr/share/locale/locale.alias ja_JP.UTF-8

# Go, gooseをインストール
RUN apt-get update \
  && apt-get install -y wget \
  && wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz \
  && tar -xvf go1.21.0.linux-amd64.tar.gz \
  && mv go /usr/local \
  && /usr/local/go/bin/go install github.com/pressly/goose/v3/cmd/goose@latest

ENV LANG ja_JP.utf8
ENV PATH $PATH:/usr/local/go/bin:/root/go/bin