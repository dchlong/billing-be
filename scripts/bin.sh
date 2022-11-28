#!/usr/bin/env bash

SCRIPTPATH="$(
    cd "$(dirname "$0")"
    pwd -P
)"

CURRENT_DIR=$SCRIPTPATH
ROOT_DIR="$(dirname $CURRENT_DIR)"
export GOBIN=$ROOT_DIR/bin
export PATH=$PATH:$GOBIN

function install_mockgen() {
  command -v $GOBIN/mockgen >/dev/null 2>&1 || {
    echo ""
    echo "tkit-cli is installing mockgen@v1.6.0"
    go install github.com/golang/mock/mockgen@v1.6.0
  }
}

function install_wire() {
  command -v $GOBIN/wire >/dev/null 2>&1 || {
    echo ""
    echo "tkit-cli is installing wire@v0.5.0"
    go install github.com/google/wire/cmd/wire@v0.5.0
  }
}

function install_golangci_lint() {
  command -v $GOBIN/golangci-lint >/dev/null 2>&1 || {
    echo ""
    echo "tkit-cli is installing golangci-lint@v1.50.0"
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1
  }
}

function dc_infra() {
    COMPOSE_FILE=$ROOT_DIR/deployments/docker-compose.test.yml
    docker-compose -p billing -f $COMPOSE_FILE $@ || exit 1
}

function run_integration_test() {
  echo "Starting infrastructure..."
  dc_infra up -d
  echo 'Running integration test'
  export CONFIG_FILE=$ROOT_DIR/configs/app.yaml
  export $(grep -v '^#' "$ROOT_DIR/deployments/base.env" | xargs -0) >/dev/null 2>&1
  go test -p 1 ./integration-test/... || {
      echo 'testing failed'
      exit 1
  }

  dc_infra down
}

function run_unit_test() {
  echo 'Running unit test'
  # note we use -p 1 to make sure that we only test 1 package at the same time
  # if we test >= 2 packages at same time
  go test -p 1 ./internal/... || {
      echo 'testing failed'
      exit 1
  }
}

function setup_env_variables() {
  export $(grep -v '^#' "$ROOT_DIR/deployments/base.env" | xargs -0) >/dev/null 2>&1
}

function run_integration_test() {
  echo "Starting infrastructure..."
  dc_infra up -d
  echo 'Running integration test'
  export CONFIG_FILE=$ROOT_DIR/configs/app.yaml
  setup_env_variables
  go test -p 1 ./integration-test/... || {
      echo 'testing failed'
      exit 1
  }

  dc_infra down
}


function run_docker() {
  COMPOSE_FILE=$ROOT_DIR/deployments/docker-compose.local.yml
  docker-compose -f $COMPOSE_FILE build
  docker-compose -p billing_local -f $COMPOSE_FILE up $@
}

function run() {
  case $1 in
    docker)
      run_docker ${@:2}
      ;;
    *)
      setup_env_variables
      go run cmd/app/main.go
      ;;
  esac

}

function generate() {
    install_mockgen
    install_wire
    go generate ./...
}

function lint() {
    install_golangci_lint
    ./bin/golangci-lint run --timeout 10m ./...
}

case $1 in
generate)
    generate
    ;;
unit_test)
    run_unit_test
    ;;
integration_test)
    run_integration_test
    ;;
test)
    run_unit_test
    run_integration_test
    ;;
lint)
    lint
    ;;
run)
    run ${@:2}
    ;;
*)
    echo "./scripts/bin.sh [generate|unit_test|integration_test|test|lint|run]"
    ;;
esac
