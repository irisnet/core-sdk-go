docker run -d -it --rm -p 36657:26657 -p 9190:9090 --name coresdktest1 coresdktest

docker exec -i -d coresdktest1 /bin/sh -c "iris start --home /scripts/test1/node0/iris"

sleep 5

rm -rf ts-relayer
ibc-setup init --home ts-relayer
echo "------------------------------ setup init success ------------------------------"
echo "
version: 1

chains:
  test:
    chain_id: test
    prefix: iaa
    gas_price: 0.1uiris
    hd_path: m/44'/118'/0'/0/0
    ics20_port: 'transfer'
    rpc:
      - http://127.0.0.1:26657
  test1:
    chain_id: test1
    prefix: iaa
    gas_price: 0.1uiris
    hd_path: m/44'/118'/0'/0/0
    ics20_port: 'transfer'
    rpc:
      - http://127.0.0.1:36657" > ts-relayer/registry.yaml

echo "------------------------------ registry.yaml success ------------------------------ "

address=$(ibc-setup init --home=ts-relayer --src=test --dest=test1)
address=${address:123}
echo "------------------------------ setup init key success ,\n address:${address} ------------------------------"

#send some token to relayer
echo '1234567890\n'|iris tx bank send iaa1x6nrhlx2he9kw73x8qcwgsl9tznh6z52msdkwx  ${address}  100000iris --keyring-backend file --home test/node0/iriscli --chain-id=test --broadcast-mode=block -y --node http://localhost:26657
echo '1234567890\n'|iris tx bank send iaa10njupdhmnyma2s7ghcapgtnw9kzg9gkjdylyla  ${address}  100000iris --keyring-backend file --home test1/node0/iriscli --chain-id=test1 --broadcast-mode=block -y  --node http://localhost:36657
echo "------------------------------ send money success ------------------------------"
ibc-setup ics20 --home=ts-relayer
echo "setup connection success"
echo '1234567890\n'|ibc-relayer start --home=ts-relayer --src=test --dest=test1 --src-connection=connection-0 --dest-connection=connection-0 --enable-metrics  -v --poll=10 --max-age-src=200 --max-age-dest=200 >  ts-relayer/relayer.log 2>&1 &
echo "relayer start ！ log in ts-relayer/relayer.log \n use \" tail -f ts-relayer/relayer.log \" to look"
# echo '1234567890\n'|iris tx ibc-transfer transfer transfer channel-0 iaa10njupdhmnyma2s7ghcapgtnw9kzg9gkjdylyla 10000000uatom --home test/node0/iris --keyring-dir test/node0/iriscli --keyring-backend=file --node=tcp://localhost:26657 --from=validator --chain-id=test -b block -y
