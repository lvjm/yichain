# pull images from docker hub and push them to ali docker registry
# this script will save hours for fabric environment setup on mainland servers
# run it from an oversea server
#!/bin/bash
set -x

# pull latest images
VERSION=1.3.0
MARCH=amd64
TAG=$MARCH-$VERSION
THIRDPARTY_VERSION=0.4.13
THIRDPARTY_TAG=$MARCH-$THIRDPARTY_VERSION

for IMAGES in peer orderer ccenv javaenv tools ca; do
    echo "PULLING FABRIC IMAGE: $IMAGES"
    echo
    docker pull hyperledger/fabric-$IMAGES:$TAG
    #docker tag hyperledger/fabric-$IMAGES:$TAG hyperledger/fabric-$IMAGES:latest
done

for IMAGES in couchdb kafka zookeeper; do
    echo "PULLING FABRIC IMAGE: $IMAGES"
    echo
    docker pull hyperledger/fabric-$IMAGES:$THIRDPARTY_TAG
    #docker tag hyperledger/fabric-$IMAGES:$VERSION hyperledger/fabric-$IMAGES:latest
done

# push images to ali docker registry
docker login registry.cn-hangzhou.aliyuncs.com

for IMAGES in peer orderer ccenv javaenv tools ca; do
    echo "==> PUSHING FABRIC IMAGE: $IMAGES"
    echo
    docker tag hyperledger/fabric-$IMAGES:$TAG registry.cn-hangzhou.aliyuncs.com/at2chain/fabric-$IMAGES:$TAG
    docker push registry.cn-hangzhou.aliyuncs.com/at2chain/fabric-$IMAGES:$TAG
    #docker tag hyperledger/fabric-$IMAGES:$TAG registry.cn-hangzhou.aliyuncs.com/at2chain/fabric-$IMAGES:latest
    #docker push registry.cn-hangzhou.aliyuncs.com/at2chain/fabric-$IMAGES:latest
done

for IMAGES in couchdb kafka zookeeper; do
    echo "==> PUSHING FABRIC IMAGE: $IMAGES"
    echo
    docker tag hyperledger/fabric-$IMAGES:$THIRDPARTY_TAG registry.cn-hangzhou.aliyuncs.com/at2chain/fabric-$IMAGES:$THIRDPARTY_TAG
    docker push registry.cn-hangzhou.aliyuncs.com/at2chain/fabric-$IMAGES:$THIRDPARTY_TAG
    #docker tag hyperledger/fabric-$IMAGES:$THIRDPARTY_TAG registry.cn-hangzhou.aliyuncs.com/at2chain/fabric-$IMAGES:latest
    #docker push registry.cn-hangzhou.aliyuncs.com/at2chain/fabric-$IMAGES:latest
done
