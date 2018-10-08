#!/bin/sh
while true; do
  ./mirafeed
  git add .
  git diff-index --quiet HEAD -- || git commit --no-gpg-sign -m 'Feed fixed' && git push
  result="$(./feedvalidator/src/demo.py https://ameer.io/test/FixedFeed.rss | grep -v Validating | grep -v guid)"
  if [ "$result" != "" ]; then
    echo "Errors detected in feed!"
    curl --data-urlencode "value1=$result" https://maker.ifttt.com/trigger/mirafeed_error/with/key/hAVKuiLxTZFFNyiGtd1FubyVsOwTOHzWzJocBA0dCJs
  fi
  echo ""
  echo "Sleeping for half an hour"
  sleep 30m
done
