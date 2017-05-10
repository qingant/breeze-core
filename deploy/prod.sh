#!/usr/bin/env sh
set -x
docker stack rm breeze
sleep 5
docker stack deploy -c docker-compose.yml breeze  --with-registry-auth
docker service update --force breeze_python
docker service update --force breeze_core

docker service update --publish-rm 5000:5000 breeze_python
docker service update --endpoint-mode dnsrr breeze_python
docker service update --host-add "tsdb.stockpalm.com:10.10.181.157" breeze_python
