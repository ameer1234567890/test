#!/bin/sh
while true; do
  ./mirafeed
  git add .
  git diff-index --quiet HEAD -- || git commit --no-gpg-sign -m 'Feed fixed' && git push
  sleep 10
done
