#!/bin/bash

function get_image_tag() {
    local currHash=$(git log -1 --pretty=format:'%h' 2>/dev/null)
    local currBranch=$(git branch 2>/dev/null | grep ^\* | awk '{print $2}')

    if [[ -z $currBranch ]] || [[ -z $currHash ]]; then
            echo "current branch or current commit id is empty"
            exit
    fi
    echo -n "${currBranch}-${currHash}"
}

Program="$0"
HomeDir=$(cd `dirname $Program`; pwd)

cd "$HomeDir"

usage() {
    cat << EOF
$Program run|test [-e xx=yy[ -e aa=bb[ ...]]] [-v[ -count=1]]
-e xx=yy -e aa=bb will set env

Examples:
    ./build.sh test -v -count=1 -e redis_host_port=127.0.0.1:6379 -e redis_db=1 -e redis_pw=123456
    ./build.sh run  -e redis_host_port=127.0.0.1:6379 -e redis_db=1 -e redis_pw=123456
EOF
}

usage_exit() {
    usage
    exit 1
}

if [[ $# -eq 0 ]]; then
    usage_exit
fi

mode=${1:-build}
shift
echo "Enter mode: $mode"

goOptions=()

# 设置环境变量
while [[ $# -gt 0 ]]; do
    if [[ $1 == '-e' ]]; then
        key=${2%%=*}
        value=${2##*=}
        export $key=$value
        shift 2
    else
        goOptions+=($1)
        shift
    fi
done

case "$mode" in
    run)
        go run "${goOptions[@]}" "$HomeDir/main.go"
    ;;
    test)
        go test "${goOptions[@]}" "$HomeDir"/...
    ;;
    build)
        go build -o captcha "${goOptions[@]}" "$HomeDir"/main.go
        image="leeoj2/captcha:`get_image_tag`"
        docker build -t "$image" .
        docker push "$image"
    ;;
    *)
        usage_exit
    ;;
esac
