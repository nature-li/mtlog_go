export GOPATH=$(pwd $(dirname $0))
echo ${GOPATH}
go install main
