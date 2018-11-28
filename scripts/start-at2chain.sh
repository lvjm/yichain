#!/bin/bash
set -x
#trap read debug

export FABRIC_CFG_PATH=${PWD}
#export FABRIC_CFG_PATH=$(dirname `pwd`)

mkdir -p channel-artifacts

cryptogen generate --config=crypto-config.yaml

configtxgen -profile At2ChainOrdererGenesis -outputBlock ./channel-artifacts/genesis.block

configtxgen -profile TraceChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID tracec

configtxgen -profile TraceChannel -outputAnchorPeersUpdate ./channel-artifacts/CoreMSPanchors.tx -channelID tracec -asOrg CoreMSP

configtxgen -profile TraceChannel -outputAnchorPeersUpdate ./channel-artifacts/SupportMSPanchors.tx -channelID tracec -asOrg SupportMSP

#different sed flag between linux and macos
ARCH=`uname -s | grep Darwin`
if [ -z $ARCH ]; then
  OPTS="-i"
elif [ $ARCH == 'Darwin' ]; then
  OPTS="-it"
else
  OPTS="-i"
fi
#cd crypto-config/peerOrganizations/core.at2chain.com/ca
CURRENT_DIR=$PWD
cd crypto-config/peerOrganizations/core.at2chain.com/ca/
PRIV_KEY=$(ls *_sk)
cd "$CURRENT_DIR"
sed $OPTS "s/CA1_PRIVATE_KEY/${PRIV_KEY}/g" docker-compose.yaml
sed $OPTS "s/CA1_PRIVATE_KEY/${PRIV_KEY}/g" docker-compose-with-explorer.yaml

set +x
