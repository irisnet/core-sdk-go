package types

import (
	"crypto/x509"
	"fmt"
	"os"

	"github.com/irisnet/core-sdk-go/crypto/keyring"

	store2 "github.com/irisnet/core-sdk-go/store"

	"github.com/cosmos/cosmos-sdk/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	defaultGas           = 200000
	defaultFees          = "4iris"
	defaultTimeout       = 5
	defaultLevel         = "info"
	defaultMaxTxsBytes   = 1073741824
	defaultAlgo          = "secp256k1"
	defaultMode          = Sync
	defaultPath          = "$HOME/irishub-sdk-go/leveldb"
	defaultGasAdjustment = 1.0
	defaultTxSizeLimit   = 1048576
	BIP44Prefix          = "44'/118'/"
	PartialPath          = "0'/0/0"
	FullPath             = "m/" + BIP44Prefix + PartialPath
)

type ClientConfig struct {
	// RPCAddr node rpc address
	RPCAddr string

	// irishub grpc address
	GRPCAddr string

	// grpc dial options
	GRPCOptions []grpc.DialOption

	// irishub chain-id
	ChainID string

	// max gas limit
	Gas uint64

	// Fee amount of point
	Fee types.DecCoins

	// PrivKeyArmor DAO Implements
	KeyDAO store2.KeyDAO

	// Private key generation algorithm(sm2,secp256k1)
	Algo string

	// Transaction broadcast Mode
	Mode BroadcastMode

	// Timeout for accessing the blockchain (such as query transactions, broadcast transactions, etc.)
	Timeout uint

	// log level(trace|debug|info|warn|error|fatal|panic)
	Level string

	// maximum bytes of a transaction
	MaxTxBytes uint64

	// adjustment factor to be multiplied against the estimate returned by the tx simulation;
	GasAdjustment float64

	// whether to enable caching
	Cached bool

	TokenManager TokenManager

	KeyManager keyring.KeyManager

	TxSizeLimit uint64

	// BIP44 path
	BIP44Path string

	// Header for rpc or http
	Header map[string][]string

	// WSAddr for ws or wss protocol
	WSAddr string

	// bech32 Address Prefix
	sdkcfg *types.Config

	FeeGranter types.AccAddress
	FeePayer   types.AccAddress
}

func NewClientConfig(rpcAddr, grpcAddr, chainID string, options ...Option) (ClientConfig, error) {
	cfg := ClientConfig{
		RPCAddr:  rpcAddr,
		ChainID:  chainID,
		GRPCAddr: grpcAddr,
	}
	for _, optionFn := range options {
		if err := optionFn(&cfg); err != nil {
			return ClientConfig{}, err
		}
	}

	if err := cfg.checkAndSetDefault(); err != nil {
		return ClientConfig{}, err
	}
	return cfg, nil
}

func (cfg *ClientConfig) checkAndSetDefault() error {
	if len(cfg.RPCAddr) == 0 {
		return fmt.Errorf("nodeURI is required")
	}

	if len(cfg.ChainID) == 0 {
		return fmt.Errorf("chainID is required")
	}

	if err := GasOption(cfg.Gas)(cfg); err != nil {
		return err
	}

	if err := FeeOption(cfg.Fee)(cfg); err != nil {
		return err
	}

	if err := AlgoOption(cfg.Algo)(cfg); err != nil {
		return err
	}

	if err := KeyDAOOption(cfg.KeyDAO)(cfg); err != nil {
		return err
	}

	if err := ModeOption(cfg.Mode)(cfg); err != nil {
		return err
	}

	if err := TimeoutOption(cfg.Timeout)(cfg); err != nil {
		return err
	}

	if err := LevelOption(cfg.Level)(cfg); err != nil {
		return err
	}

	if err := MaxTxBytesOption(cfg.MaxTxBytes)(cfg); err != nil {
		return err
	}

	if err := TokenManagerOption(cfg.TokenManager)(cfg); err != nil {
		return err
	}

	if err := TxSizeLimitOption(cfg.TxSizeLimit)(cfg); err != nil {
		return err
	}

	if err := Bech32AddressPrefixOption(cfg.sdkcfg)(cfg); err != nil {
		return err
	}

	if err := BIP44PathOption(cfg.BIP44Path)(cfg); err != nil {
		return err
	}
	return GasAdjustmentOption(cfg.GasAdjustment)(cfg)
}

type Option func(cfg *ClientConfig) error

func FeeOption(fee types.DecCoins) Option {
	return func(cfg *ClientConfig) error {
		if fee == nil || fee.Empty() || !fee.IsValid() {
			fees, _ := types.ParseDecCoins(defaultFees)
			fee = fees
		}
		cfg.Fee = fee
		return nil
	}
}

func KeyDAOOption(dao store2.KeyDAO) Option {
	return func(cfg *ClientConfig) error {
		if dao == nil {
			defaultPath := os.ExpandEnv(defaultPath)
			levelDB, err := store2.NewLevelDB(defaultPath, nil)
			if err != nil {
				return err
			}
			dao = levelDB
		}
		cfg.KeyDAO = dao
		return nil
	}
}

func GasOption(gas uint64) Option {
	return func(cfg *ClientConfig) error {
		if gas <= 0 {
			gas = defaultGas
		}
		cfg.Gas = gas
		return nil
	}
}

func AlgoOption(algo string) Option {
	return func(cfg *ClientConfig) error {
		if algo == "" {
			algo = defaultAlgo
		}
		cfg.Algo = algo
		return nil
	}
}

func ModeOption(mode BroadcastMode) Option {
	return func(cfg *ClientConfig) error {
		if mode == "" {
			mode = defaultMode
		}
		cfg.Mode = mode
		return nil
	}
}

func TimeoutOption(timeout uint) Option {
	return func(cfg *ClientConfig) error {
		if timeout <= 0 {
			timeout = defaultTimeout
		}
		cfg.Timeout = timeout
		return nil
	}
}

func LevelOption(level string) Option {
	return func(cfg *ClientConfig) error {
		if level == "" {
			level = defaultLevel
		}
		cfg.Level = level
		return nil
	}
}

func MaxTxBytesOption(maxTxBytes uint64) Option {
	return func(cfg *ClientConfig) error {
		if maxTxBytes <= 0 {
			maxTxBytes = defaultMaxTxsBytes
		}
		cfg.MaxTxBytes = maxTxBytes
		return nil
	}
}

func GasAdjustmentOption(gasAdjustment float64) Option {
	return func(cfg *ClientConfig) error {
		if gasAdjustment <= 0 {
			gasAdjustment = defaultGasAdjustment
		}
		cfg.GasAdjustment = gasAdjustment
		return nil
	}
}

func CachedOption(enabled bool) Option {
	return func(cfg *ClientConfig) error {
		cfg.Cached = enabled
		return nil
	}
}

func TokenManagerOption(tokenManager TokenManager) Option {
	return func(cfg *ClientConfig) error {
		if tokenManager == nil {
			tokenManager = DefaultTokenManager{}
		}
		cfg.TokenManager = tokenManager
		return nil
	}
}

func TxSizeLimitOption(txSizeLimit uint64) Option {
	return func(cfg *ClientConfig) error {
		if txSizeLimit <= 0 {
			txSizeLimit = defaultTxSizeLimit
		}
		cfg.TxSizeLimit = txSizeLimit
		return nil
	}
}

func KeyManagerOption(keyManager keyring.KeyManager) Option {
	return func(cfg *ClientConfig) error {
		cfg.KeyManager = keyManager
		return nil
	}
}

func Bech32AddressPrefixOption(sdkcfg *types.Config) Option {
	return func(cfg *ClientConfig) error {
		if sdkcfg != nil {
			sdkcfg.Seal()
		}
		return nil
	}
}

func BIP44PathOption(bIP44Path string) Option {
	return func(cfg *ClientConfig) error {
		if bIP44Path == "" {
			bIP44Path = FullPath
		}
		cfg.BIP44Path = bIP44Path
		return nil
	}
}

func HeaderOption(header map[string][]string) Option {
	return func(cfg *ClientConfig) error {
		cfg.Header = header
		return nil
	}
}

func WSAddrOption(wsAddr string) Option {
	return func(cfg *ClientConfig) error {
		cfg.WSAddr = wsAddr
		return nil
	}
}

func GRPCOptions(gRPCOptions []grpc.DialOption, TLS bool, rpcAddr string) Option {
	return func(cfg *ClientConfig) error {
		if !TLS {
			cfg.GRPCOptions = gRPCOptions
			return nil
		}

		certificateList, err := GetTLSCertPool(rpcAddr)
		if err != nil {
			panic(err)
		}

		roots := x509.NewCertPool()
		for i := range certificateList {
			roots.AddCert(certificateList[i])
		}
		cert := credentials.NewClientTLSFromCert(roots, "")
		cfg.GRPCOptions = append(gRPCOptions, grpc.WithTransportCredentials(cert))

		return nil
	}
}

func FeeGranterOptions(feeGranter string) Option {
	return func(cfg *ClientConfig) error {
		granter, err := types.AccAddressFromBech32(feeGranter)
		if err != nil {
			panic(err)
		}
		cfg.FeeGranter = granter
		return nil
	}
}

func FeePayerOptions(feePayer string) Option {
	return func(cfg *ClientConfig) error {
		feePayer, err := types.AccAddressFromBech32(feePayer)
		if err != nil {
			panic(err)
		}
		cfg.FeePayer = feePayer
		return nil
	}
}
