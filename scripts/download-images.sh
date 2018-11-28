
#!/bin/bash
VERSION=1.3.0
THIRDPARTY_TAG=0.4.13
for IMAGES in peer orderer ccenv javaenv tools ca;do
    echo "==> FABRIC IMAGE: $IMAGES"
    echo
    docker pull hyperledger/fabric-$IMAGES:$VERSION
done
for IMAGES in couchdb kafka zookeeper; do
      echo "==> THIRDPARTY DOCKER IMAGE: $IMAGES"
      echo
      docker pull hyperledger/fabric-$IMAGES:$THIRDPARTY_TAG
done

