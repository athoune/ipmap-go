# IPmap

Expose https://ipmap.ripe.net/docs/01.manual/#full-dumps with https://github.com/athoune/iptree

## Test it

    make
    ./bin/ipmap geolocations_2022-04-04.csv.bz2

In another terminals :

    nc 127.0.0.1 1234

You can put IP address, and press [enter].
