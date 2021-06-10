#!/bin/bash

function usage() {
    echo "Usage:"
    echo "  ./$(basename $0) <HOST1:PORT1> <interval in seconds>" ;
    echo "E.g."
    echo "  ./$(basename $0) 10.20.30.40:80 30";
    echo "This executes curl liveness test on services on \ ";
    echo " 10.20.30.40:80 with an interval of 30 secs"
    exit -1
}

function main {
    if [ $# -ne 2 ]; then
        usage;
    fi
    host=$1
    interval=$2
    while true; do
        date
        curl -m 5 -I $host 
        sleep $interval 
    done
}

main "$@"
