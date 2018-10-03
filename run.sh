./mirafeed
git add .
git diff-index --quiet HEAD -- || git commit --no-gpg-sign -m 'Feed fixed' && git push
