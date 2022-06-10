#!/bin/sh

if [ $# -ne 1 ]; then
	echo "$0 <tag>"
	exit 0
fi

tag=$1

echo "tag: ${tag}"

docker build -t pubplatform:${tag} .

docker tag pubplatform:${tag} leeoj2/pubplatform:${tag}
docker push leeoj2/pubplatform:${tag}