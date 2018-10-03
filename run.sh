./mirafeed
git diff-index --quiet HEAD -- || git add . && git commit --no-gpg-sign -m 'Feed fixed' && git push
