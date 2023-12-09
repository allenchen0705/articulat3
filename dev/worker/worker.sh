#!/usr/bin/env bash

set -euo pipefail

unset -v username
unset -v server

username=${USERNAME}
temporal_server_host_port=${TEMPORAL_SERVER_HOST_PORT}

if [ -z "$username" ]; then
    echo "USERNAME is not set"
    exit 1
fi

if [ -z "$temporal_server_host_port" ]; then
    echo "TEMPORAL_SERVER_HOST_PORT is not set"
    exit 1
fi

require() {
  if ! hash "$1" &>/dev/null; then
    echo "'$1' not found in PATH"
    exit 1
  fi
}

require singularity
require nvidia-smi

# Remove build artifacts
rm -rf build-temp*

project_name="project-44-toronto-intelligence-m"
app_repo_dir="/home/$username/$project_name"
model_repo_dir="/home/$username/MVDream-threestudio"

ssh-keyscan github.com >> ~/.ssh/known_host

echo "Cloning projects..."
if [ ! -d $app_repo_dir ]; then
    git clone git@github.com:csc301-2023-fall/project-44-toronto-intelligence-m.git
fi

if [ ! -d $model_repo_dir ]; then
    git clone git@github.com:bytedance/MVDream-threestudio.git
fi

if [ ! -d "$model_repo_dir/extern" ]; then
    git clone https://github.com/bytedance/MVDream extern/MVDream
    # pip install -e extern/MVDream
fi

cd $model_repo_dir
git pull

# Need to patch the version of xformers to work with the driver version on
# tisl machines.
set +e
patch -r - -u -f requirements.txt <<'EOF'
--- requirements.txt	2023-12-05 22:47:45
+++ requirements.p	2023-12-05 22:48:54
@@ -25,7 +25,7 @@
 torchmetrics

 # deepfloyd
-xformers
+xformers==0.0.20
 bitsandbytes
 sentencepiece
 safetensors
EOF
set -e

cd $app_repo_dir
git pull

cd temporal


box_dir="/home/$username/box"

echo "Building singularity environment..."
if [ ! -d $box_dir ]; then
    singularity build --sandbox --fakeroot $box_dir Singularity.def
else
    echo "Singularity environment exists."
fi

echo "Starting singularity worker..."

cat <<EOT > /tmp/run-worker-process.sh
cd /root/MVDream-threestudio
pip install -e extern/MVDream

cd /root/$project_name/temporal
export TEMPORAL_SERVER_HOST_PORT=${TEMPORAL_SERVER_HOST_PORT}
/root/box/usr/local/go/bin/go run worker/main.go
EOT

chmod 700 /tmp/run-worker-process.sh

singularity exec --fakeroot --writable --nv --network "host" /home/$username/box /tmp/run-worker-process.sh

wait
