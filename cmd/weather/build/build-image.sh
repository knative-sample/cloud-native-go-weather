#!/bin/bash
#****************************************************************#
# Create Date: 2019-02-02 22:16
#********************************* ******************************#

ROOTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

NAME="appos"
TAG="v1-$(date +%Y''%m''%d''%H''%M''%S)"

docker build -t "${NAME}:${TAG}" -f ${ROOTDIR}/Dockerfile ${ROOTDIR}/../../../

##array=( registry.cn-beijing.aliyuncs.com  registry.cn-hangzhou.aliyuncs.com registry.cn-huhehaote.aliyuncs.com registry.cn-shanghai.aliyuncs.com registry.cn-shenzhen.aliyuncs.com  registry.cn-qingdao.aliyuncs.com registry.cn-zhangjiakou.aliyuncs.com registry.ap-southeast-2.aliyuncs.com )
#array=( registry.cn-hangzhou.aliyuncs.com )
#for registry in "${array[@]}"
#do
#    echo "push images to ${registry}/kstarter/${NAME}:${TAG}"
#    docker tag "${NAME}:${TAG}" "${registry}/kstarter/${NAME}:${TAG}"
#    docker push "${registry}/kstarter/${NAME}:${TAG}"
#
#    docker tag "${NAME}:${TAG}" "${registry}/kstarter/${NAME}:latest"
#    docker push "${registry}/kstarter/${NAME}:latest"
#done
