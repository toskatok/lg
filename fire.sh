#!/bin/bash
trap 'end' INT

end() {
        echo "Let's finish"

        pkill lg
}

usage() {
        echo "fire set your platform on fire :D"
        echo "$0 -b [broker] -r [rate]"
        echo "-b [broker]: broker url e.g. 127.0.0.1:1883"
        echo "-r [rate]: rate e.g. 10s"
}

main() {
        local rate
        local broker

        rate=$1
        broker=$2

        for i in $(seq 1 100); do
                dev=$((i + 10))
                ./lg --rate $rate --deveui $dev --broker $broker &
                if [ $? ]; then
                        echo "lg $i: $!"
                else
                        exit
                fi
                echo
        done

        cat
}


rate="10s"
broker="127.0.0.1:1883"

while getopts 'r:b:' argv; do
        case $argv in
                r)
                        rate=$OPTARG
                        ;;
                b)
                        broker=$OPTARG
                        ;;
        esac
done

main $rate $broker
