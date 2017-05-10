FROM debian:jessie
MAINTAINER M.T matao.xjtu@gmail.com

LABEL version="0.1"
ADD ./breeze-core /usr/local/bin/

ADD ./config.toml /opt/app/breeze/
ADD ./config.prod.toml /opt/app/breeze/
ADD ./config.dev.toml /opt/app/breeze/
ADD ./config.staging.toml /opt/app/breeze/
ADD ./examples /opt/app/breeze/examples/
WORKDIR /opt/app/breeze/

EXPOSE 8080
EXPOSE 6379
ENTRYPOINT ["breeze-core"]
