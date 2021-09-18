# Core SDK

Golang SDK for Tendermint & Cosmos-SDK Core Modules

## install

### Requirement

Go version above 1.16.4

### Use Go Mod

```go
replace (
    github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
    github.com/tendermint/tendermint => github.com/bianjieai/tendermint v0.34.1-irita-210113
)
```

### Init Client

The initialization SDK code is as follows:

```go
    bech32AddressPrefix := types.AddrPrefixCfg{
        AccountAddr:   "iaa",
        ValidatorAddr: "iva",
        ConsensusAddr: "ica",
        AccountPub:    "iap",
        ValidatorPub:  "ivp",
        ConsensusPub:  "icp",
    }
    options := []types.Option{
        types.KeyDAOOption(store.NewMemory(nil)),
        types.TimeoutOption(10),
        types.TokenManagerOption(TokenManager{}),
        types.KeyManagerOption(crypto.NewKeyManager()),
        types.Bech32AddressPrefixOption(bech32AddressPrefix),
        types.BIP44PathOption(""),
    }
    cfg, err := types.NewClientConfig(nodeURI, grpcAddr, chainID, options...)
    if err != nil {
        panic(err)
    }
    
    s.Client = sdk.NewClient(cfg)
```

The NewClientConfig component mainly contains the parameters used in the SDK, the specific meaning is shown in the table
below.

| Iterm     | Type           | Description                                                                                         |
| --------- | -------------- | --------------------------------------------------------------------------------------------------- |
| NodeURI   | string         | The RPC address of the IRIShub node connected to the SDK, for example: localhost: 26657             |
| GRPCAddr  | string         | The GRPC address of the IRIShub node connected to the SDK, for example: localhost: 9090             |
| Network   | enum           | IRIShub network type, value: Testnet, Mainnet                                                       |
| ChainID   | string         | ChainID of IRIShub, for example: IRIShub                                                            |
| Gas       | uint64         | The maximum gas to be paid for the transaction, for example: 20000                                  |
| Fee       | DecCoins       | Transaction fees to be paid for transactions                                                        |
| KeyDAO    | KeyDAO         | Private key management interface, If the user does not provide it, the default LevelDB will be used |
| Mode      | enum           | Transaction broadcast mode, value: Sync, Async, Commit                                              |
| StoreType | enum           | Private key storage method, value: Keystore, PrivKey                                                |
| Timeout   | time. Duration | Transaction timeout, for example: 5s                                                                |
| Level     | string         | Log output level, for example: info                                                                 |

If you want to use SDK to send a transfer transaction, the example is as follows:

There is more example of query and send tx

```go
    coins, err := types.ParseDecCoins("10iris")
    to := "iaa1hp29kuh22vpjjlnctmyml5s75evsnsd8r4x0mm"
    baseTx := types.BaseTx{
        From:               s.Account().Name,
        Gas:                200000,
        Memo:               "TEST",
        Mode:               types.Commit,
        Password:          "password",
        SimulateAndExecute: false,
        GasAdjustment:      1.5,
    }
    
    res, err := s.Bank.Send(to, coins, baseTx)
```
