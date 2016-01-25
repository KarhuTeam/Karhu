#!/bin/bash

KARHU_DATA_DIR=${DATA_DIR:-/data}

echo "Ensure data dir exist"
mkdir -p $KARHU_DATA_DIR


echo "DEBUG=${DEBUG:-1}
MGO_HOSTS=${MGO_HOSTS:-mongo}
MGO_DB=${MGO_DB:-karhu}

PUBLIC_HOST=${PUBLIC_HOST:-http://127.0.0.1:8080}

STORAGE_DRIVER=${STORAGE_DRIVER:-fs}
DATA_DIR=$KARHU_DATA_DIR" > .env

cat .env; echo

./karhu
