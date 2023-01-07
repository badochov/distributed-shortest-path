worker

if [ -z $1 ]; then
  echo "Usage $0 <version>"
  exit 1
fi

VERSION=$1

for i in {0..15}; do
  kubectl set image --namespace shortest-path "deployment/workers-region-$i" "worker=shortest-path/worker:$VERSION"
done