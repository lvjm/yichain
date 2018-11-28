# this script will download fabric binaries to go bin
#!/bin/bash
set -x

mkdir -p $HOME/go

VERSION=1.3.0
ARCH=$(echo "$(uname -s|tr '[:upper:]' '[:lower:]'|sed 's/mingw64_nt.*/windows/')-$(uname -m | sed 's/x86_64/amd64/g')" | awk '{print tolower($0)}')
#Set MARCH variable i.e ppc64le,s390x,x86_64,i386
MARCH=`uname -m`

echo "===> Downloading platform binaries"
curl -OL https://nexus.hyperledger.org/content/repositories/releases/org/hyperledger/fabric/hyperledger-fabric/${ARCH}-${VERSION}/hyperledger-fabric-${ARCH}-${VERSION}.tar.gz
tar -xzf hyperledger-fabric-${ARCH}-${VERSION}.tar.gz -C $HOME/go/
#rm hyperledger-fabric-${ARCH}-${VERSION}.tar.gz

set +x
