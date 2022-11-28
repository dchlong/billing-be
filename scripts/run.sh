#!/usr/bin/env sh

case $1 in
app)
  chmod +x /source/app
  /source/app
  ;;
*)
  echo "./scripts/run.sh [app]"
  ;;
esac
