#!/bin/sh
while true; do
  ./mirafeed
  git add .
  git diff-index --quiet HEAD -- || git commit --no-gpg-sign -m 'Feed fixed' && git push
  echo "Sleeping for half an hour"
  sleep 10
done
