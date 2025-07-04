#!/bin/bash
set -e

ACCOUNT_ID=102800182496
REGION=us-east-1
REPO_URI="${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com/particles"

get_git_hash() {
    if git rev-parse --git-dir > /dev/null 2>&1; then
        git rev-parse --short HEAD
    else
        log "⚠️" "Not a git repository, using 'latest' as version"
        echo "latest"
    fi
}

GIT_COMMIT=$(get_git_hash)
aws ecr get-login-password --region $REGION | docker login --username AWS --password-stdin $REPO_URI
docker buildx build --platform linux/amd64,linux/arm64 -t $REPO_URI:latest -t $REPO_URI:$GIT_COMMIT --push .

cd aws
cdk deploy ParticlesStack --parameters ImageTag=$GIT_COMMIT
echo "$GIT_COMMIT Deployed!"
