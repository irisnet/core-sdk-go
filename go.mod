module github.com/irisnet/core-sdk-go

go 1.16

require (
	github.com/99designs/keyring v1.1.6
	github.com/DataDog/zstd v1.4.8 // indirect
	github.com/StackExchange/wmi v1.2.1 // indirect
	github.com/armon/go-metrics v0.3.9
	github.com/avast/retry-go v3.0.0+incompatible
	github.com/bluele/gcache v0.0.2
	github.com/btcsuite/btcd v0.22.0-beta
	github.com/btcsuite/btcutil v1.0.3-0.20201208143702-a53e38424cce
	github.com/cosmos/cosmos-sdk v0.44.3
	github.com/cosmos/go-bip39 v1.0.0
	github.com/dvsekhvalnov/jose2go v0.0.0-20201001154944-b09cfaf05951
	github.com/ethereum/go-ethereum v1.10.16
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.5.2
	github.com/google/btree v1.0.1 // indirect
	github.com/gorilla/websocket v1.4.2
	github.com/magiconair/properties v1.8.5
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mtibben/percent v0.2.1
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.17.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/common v0.29.0
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475
	github.com/regen-network/cosmos-proto v0.3.1
	github.com/rs/cors v1.8.0 // indirect
	github.com/sasha-s/go-deadlock v0.2.1-0.20190427202633-1595213edefa // indirect
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/viper v1.9.0 // indirect
	github.com/stretchr/testify v1.8.1
	github.com/tendermint/crypto v0.0.0-20191022145703-50d29ede1e15
	github.com/tendermint/go-amino v0.16.0
	github.com/tendermint/tendermint v0.34.14
	github.com/tendermint/tm-db v0.6.4
	github.com/tjfoc/gmsm v1.4.0
	github.com/tklauser/go-sysconf v0.3.7 // indirect
	go.etcd.io/bbolt v1.3.6 // indirect
	golang.org/x/crypto v0.1.0
	golang.org/x/net v0.6.0 // indirect
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f
	google.golang.org/grpc v1.53.0
	google.golang.org/protobuf v1.28.1
	gopkg.in/yaml.v2 v2.4.0
)

replace (
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
	github.com/tendermint/tendermint => github.com/bianjieai/tendermint v0.34.1-irita-210113
	github.com/tharsis/ethermint v0.8.1 => github.com/bianjieai/ethermint v0.8.2-0.20220211020007-9ec25dde74d4
)
