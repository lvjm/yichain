#!/bin/bash
set -x

cd $HOME/go/src/bitbucket.org
rm -rf at2chain
git clone --depth=1 https://bitbucket.org/AI-CORE/at2chain.git
rm -rf at2chain/channel-artifacts/*
rm -rf at2chain/crypto-config/*
cp at2chain/docker-compose-template.yaml at2chain/docker-compose.yaml
chmod 777 at2chain
