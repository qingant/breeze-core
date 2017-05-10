#!/usr/bin/bash
set -x
rm -f dist.zip
zip -o dist.zip -r * -X deploy.sh
curl -v -XPOST  -H "Content-Type: application/zip" gw.stockpalm.com:9091/api/v1/deploy --data-binary @dist.zip
