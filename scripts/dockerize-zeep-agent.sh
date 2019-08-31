#!/usr/bin/env bash

set -euo pipefail

#
# VARIABLES =====================================
#

# project root path
workspace="$(dirname $0)/.."

# target image name:tag
img_ref="mee6aas/zeep:latest"

# dockerfile path to build image
dockerfile="$workspace/Dockerfile"



#
# CHECK PRECONDITION ============================
#

# check if Dockerfile exists
if [ ! -f  $dockerfile ]; then
    echo "Dockerfile is not provided"
    exit 1
fi



#
# MAIN ==========================================
#
docker build \
    --tag  $img_ref \
    --file $dockerfile \
    $workspace
