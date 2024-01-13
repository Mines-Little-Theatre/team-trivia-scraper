#!/bin/sh
set -x
TMPDIR=$(mktemp -d)

GOOS=linux GOARCH=arm64 go build -v -trimpath -ldflags "-s" -o $TMPDIR/bootstrap github.com/Mines-Little-Theatre/team-trivia-scraper/aws-lambda && \
zip -j $TMPDIR/deploy.zip $TMPDIR/bootstrap && \
aws lambda update-function-code --function-name mlt-trivia-bot --zip-file fileb://$TMPDIR/deploy.zip --no-cli-pager

rm -Rf $TMPDIR
