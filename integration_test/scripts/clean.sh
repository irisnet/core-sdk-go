docker stop coresdktest
docker stop coresdktest1
docker rmi coresdktest

rm -rf ts-relayer
port=$(ps -a | grep "ibc-relayer start")
kill ${port:0:5}