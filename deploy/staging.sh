#!/usr/bin/env sh
set -x
docker stack rm breeze_staging
sleep 5
docker stack deploy -c docker-compose.staging.yml breeze_staging  --with-registry-auth
# docker service update --publish-rm 5000:5000 breeze_python
docker service update --force breeze_staging_python
docker service update --force breeze_staging_core

docker service update --endpoint-mode dnsrr breeze_python
docker service update --host-add "tsdb.stockpalm.com:10.10.181.157" breeze_python
