module github.com/irisnet/core-sdk-go

go 1.16

require (
	github.com/avast/retry-go v3.0.0+incompatible
	github.com/bluele/gcache v0.0.2
	github.com/btcsuite/btcd v0.22.1
	github.com/btcsuite/btcutil v1.0.3-0.20201208143702-a53e38424cce
	github.com/cosmos/go-bip39 v1.0.0
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.5.2
	github.com/magiconair/properties v1.8.6
	github.com/pkg/errors v0.9.1
	github.com/prometheus/common v0.37.0 // indirect
	github.com/regen-network/cosmos-proto v0.3.1
	github.com/sirupsen/logrus v1.9.0
	github.com/stretchr/testify v1.8.0
	github.com/tendermint/crypto v0.0.0-20191022145703-50d29ede1e15
	github.com/tendermint/go-amino v0.16.0
	github.com/tendermint/tendermint v0.34.21
	github.com/tendermint/tm-db v0.6.6
	github.com/tjfoc/gmsm v1.4.0
	golang.org/x/crypto v0.0.0-20220525230936-793ad666bf5e
	google.golang.org/genproto v0.0.0-20220725144611-272f38e5d71b
	google.golang.org/grpc v1.48.0
	google.golang.org/protobuf v1.28.0
	gopkg.in/yaml.v2 v2.4.0
)

replace (
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
	github.com/tendermint/tendermint => github.com/bianjieai/tendermint v0.34.21-irita-220906
)
