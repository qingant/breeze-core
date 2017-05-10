#!/usr/bin/env sh
docker stack rm breeze
sleep 5
docker stack deploy -c docker-compose.yml breeze  --with-registry-auth
# docker service update --endpoint-mode dnsrr breeze_python

docker service update --force breeze_python
docker service update --force breeze_core

