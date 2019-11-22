#!/bin/bash
#****************************************************************#
# Create Date: 2019-02-02 22:16
#********************************* ******************************#
ROOTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

NAME="${1}"
[ -z "${1}" ] && echo "component is empty" && exit 1

NS="knative-sample"
GIT_COMMIT="$(git rev-parse --verify HEAD)"
GIT_BRANCH=`git branch | grep \* | cut -d ' ' -f2`
TAG="$(date +%Y''%m''%d''%H''%M''%S)"
if [ ! -z "${GIT_COMMIT}" -a "${GIT_BRANCH}" != " " ]; then
  TAG="${GIT_BRANCH}_${GIT_COMMIT:0:8}-$(date +%Y''%m''%d''%H''%M''%S)"
fi

docker build -t "weather-${NAME}:${TAG}" -f ${ROOTDIR}/Dockerfile-${NAME} ${ROOTDIR}/

array=( registry.cn-hangzhou.aliyuncs.com )
for registry in "${array[@]}"
do
    echo "push images to ${registry}/${NS}/weather-${NAME}:${TAG}"
    docker tag "weather-${NAME}:${TAG}" "${registry}/${NS}/weather-${NAME}:${TAG}"
    docker push "${registry}/${NS}/weather-${NAME}:${TAG}"

    docker tag "weather-${NAME}:${TAG}" "${registry}/${NS}/weather-${NAME}:latest"
    docker push "${registry}/${NS}/weather-${NAME}:latest"
done
