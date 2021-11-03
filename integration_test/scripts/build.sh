docker build -t coresdktest -f ./Dockerfile .

docker run -d -it --rm -p 26657:26657 -p 9090:9090 --name coresdktest coresdktest
docker exec -i -d coresdktest /bin/sh -c "iris start --home /scripts/test/node0/iris"
