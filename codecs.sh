#!/bin/bash

host="192.168.0.200:8080"
project="5b615954e40b0a00081e16c6"
id=1
IFS=
code="
from codec import Codec
import cbor


class ISRC(Codec):
    thing_location = 'loc'

    def decode(self, data):
        d = cbor.loads(data)

        if 'lat' in d and 'lng' in d:
            d['loc'] = self.create_location(d['lat'], d['lng'])
            del d['lat']
            del d['lng']

        return d

    def encode(self, data):
        return cbor.dumps(data)
"

echo $code

data=$(cat <<-EOF
{
        "id": "$(printf %016x $id)",
        "code": "$code"
}
EOF
)
echo $data

curl --header "Content-Type: application/json" \
        --request POST \
        --data "$data" \
        "$host/api/runners/$project/codec"
