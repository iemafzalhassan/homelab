#!/usr/bin/env bash

set -e

REGISTRY_HOST=${REGISTRY_HOST:-"localhost:5000"}
AUTH_TOKEN=${AUTH_TOKEN:-""}

if [ -z "$AUTH_TOKEN" ]; then
  echo "Testing WITHOUT auth (Dev mode)..."
else
  echo "Testing WITH auth..."
  # If auth token is provided, you would typically `docker login` here,
  # but for a basic script test using curl we can just pass the header.
fi

echo "--- Generating a dummy image ---"
mkdir -p /tmp/kyros-test
cat <<EOF > /tmp/kyros-test/Dockerfile
FROM scratch
COPY hello.txt /
CMD ["/hello"]
EOF
echo "Hello Kyros" > /tmp/kyros-test/hello.txt

echo "--- Building dummy image ---"
docker build -t kyros-test:latest /tmp/kyros-test

echo "--- Tagging for Registry ---"
docker tag kyros-test:latest ${REGISTRY_HOST}/test-org/kyros-test:latest

echo "--- Pushing to Registry ---"
# Note: For this to work with auth, docker daemon needs to be authenticated 
# or the proxy configured to accept no-auth in dev.
docker push ${REGISTRY_HOST}/test-org/kyros-test:latest

echo "--- Pulling from Registry ---"
docker rmi ${REGISTRY_HOST}/test-org/kyros-test:latest
docker pull ${REGISTRY_HOST}/test-org/kyros-test:latest

echo "--- Cleanup ---"
rm -rf /tmp/kyros-test
docker rmi kyros-test:latest
docker rmi ${REGISTRY_HOST}/test-org/kyros-test:latest

echo "--- End to end test complete! ---"
