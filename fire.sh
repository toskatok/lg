#!/bin/bash
trap 'end' INT

end() {
        echo "Let's finish"

        pkill mqttlg
}

for i in `seq 1 100`; do
        dev=$(($i + 10))
        ./mqttlg --rate 10s --deveui $dev --broker 192.168.0.200:1883 &
        echo "mqttlg $i: $!"
done

cat
