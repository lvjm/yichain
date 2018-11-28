#!bin/bash
VERSION=1.3.0
MARCH=amd64
TAG=$MARCH-$VERSION
THIRDPARTY_VERSION=0.4.13
THIRDPARTY_TAG=$MARCH-$THIRDPARTY_VERSION
for IMAGES in peer orderer ccenv javaenv tools ca;do
    echo "==> FABRIC IMAGE: $IMAGES"
    echo
    docker pull registry.cn-hangzhou.aliyuncs.com/at2chain/fabric-$IMAGES:$TAG
done
for IMAGES in couchdb kafka zookeeper; do
      echo "==> THIRDPARTY DOCKER IMAGE: $IMAGES"
      echo
      docker pull registry.cn-hangzhou.aliyuncs.com/at2chain/fabric-$IMAGES:$THIRDPARTY_TAG
done
