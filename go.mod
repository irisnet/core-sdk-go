module github.com/irisnet/core-sdk-go

go 1.16

require (
	github.com/99designs/keyring v1.1.6
	github.com/avast/retry-go v3.0.0+incompatible
	github.com/bluele/gcache v0.0.2
	github.com/btcsuite/btcd v0.21.0-beta
	github.com/cosmos/go-bip39 v1.0.0
	github.com/dvsekhvalnov/jose2go v0.0.0-20201001154944-b09cfaf05951
	github.com/enigmampc/btcutil v1.0.3-0.20200723161021-e2fb6adb2a25
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.5.2
	github.com/hdevalence/ed25519consensus v0.0.0-20210430192048-0962ce16b305
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mtibben/percent v0.2.1
	github.com/pkg/errors v0.9.1
	github.com/prometheus/common v0.14.0
	github.com/regen-network/cosmos-proto v0.3.1
	github.com/sirupsen/logrus v1.6.0
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/crypto v0.0.0-20191022145703-50d29ede1e15
	github.com/tendermint/go-amino v0.16.0
	github.com/tendermint/tendermint v0.34.12
	github.com/tjfoc/gmsm v1.4.0
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
	google.golang.org/genproto v0.0.0-20210114201628-6edceaf6022f
	google.golang.org/grpc v1.37.0
	google.golang.org/protobuf v1.26.0
	gopkg.in/yaml.v2 v2.3.0

)

replace (
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	github.com/tendermint/tendermint => github.com/bianjieai/tendermint v0.34.8-irita-210413.0.20210908054213-781a5fed16d6
)
