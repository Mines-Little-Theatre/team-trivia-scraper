#!/bin/sh
set -x
TMPDIR=$(mktemp -d)

GOOS=linux GOARCH=arm64 go build -trimpath -ldflags "-s" -o $TMPDIR/bootstrap github.com/Mines-Little-Theatre/team-trivia-scraper/aws-lambda && \
# aarch64-linux-gnu-strip $TMPDIR/bootstrap && \
zip -j $TMPDIR/deploy.zip $TMPDIR/bootstrap && \
aws lambda update-function-code --function-name $FUNCTION_NAME --zip-file fileb://$TMPDIR/deploy.zip --no-cli-pager

rm -Rf $TMPDIR
