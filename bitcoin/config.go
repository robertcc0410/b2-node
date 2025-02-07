package bitcoin

import (
	"os"
	"path"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/spf13/viper"
)

const (
	// MAINNET...
	MAINNET = "mainnet"
	// TESTNET...
	TESTNET = "testnet"
	// SIGNET...
	SIGNET = "signet"
	// SIMNET...
	SIMNET = "simnet"
	// REGTEST...
	REGTEST = "regtest"
)

// BitconConfig defines the bitcoin config
// TODO: defined different config group eg: bitcoin, bridge, indexer, commiter
type BitconConfig struct {
	// NetworkName defines the bitcoin network name
	NetworkName string `mapstructure:"network-name"`
	// RPCHost defines the bitcoin rpc host
	RPCHost string `mapstructure:"rpc-host"`
	// RPCPort defines the bitcoin rpc port
	RPCPort string `mapstructure:"rpc-port"`
	// RPCUser defines the bitcoin rpc user
	RPCUser string `mapstructure:"rpc-user"`
	// RPCPass defines the bitcoin rpc password
	RPCPass string `mapstructure:"rpc-pass"`
	// WalletName defines the bitcoin wallet name
	WalletName string `mapstructure:"wallet-name"`
	// Destination defines the taproot transaction destination address
	Destination string `mapstructure:"destination"`
	// EnableIndexer defines whether to enable the indexer
	EnableIndexer bool `mapstructure:"enable-indexer"`
	// EnableCommitter defines whether to enable the committer
	EnableCommitter bool `mapstructure:"enable-committer"`
	// IndexerListenAddress defines the address to listen on
	IndexerListenAddress string `mapstructure:"indexer-listen-address"`
	// Bridge defines the bridge config
	Bridge BridgeConfig `mapstructure:"bridge"`
	// Dsn defines the state db dsn
	StateConfig StateConfig `mapstructure:"state"`
	// Fee defines the bitcoin tx fee
	Fee int64 `mapstructure:"fee"`
	// Evm defines the evm config
	Evm EvmConfig `mapstructure:"evm"`
}

type BridgeConfig struct {
	// EthRPCURL defines the ethereum rpc url
	EthRPCURL string `mapstructure:"eth-rpc-url"`
	// EthPrivKey defines the invoke ethereum private key
	EthPrivKey string `mapstructure:"eth-priv-key"`
	// ContractAddress defines the l1 -> l2 bridge contract address
	ContractAddress string `mapstructure:"contract-address"`
	// ABI defines the l1 -> l2 bridge contract abi
	ABI string `mapstructure:"abi"`
	// GasLimit defines the  contract gas limit
	GasLimit uint64 `mapstructure:"gas-limit"`
	// AASCARegistry defines the  contract AASCARegistry address
	AASCARegistry string `mapstructure:"aa-sca-registry"`
	// AAKernelFactory defines the  contract AAKernelFactory address
	AAKernelFactory string `mapstructure:"aa-kernel-factory"`
}

type StateConfig struct {
	Host   string `mapstructure:"host"`
	Port   int    `mapstructure:"port"`
	User   string `mapstructure:"user"`
	Pass   string `mapstructure:"pass"`
	DBName string `mapstructure:"db-name"`
}

type EvmConfig struct {
	// EnableListener defines whether to enable the listener
	EnableListener bool `mapstructure:"enable-listener"`
	// Deposit defines the deposit event hash
	Deposit string `mapstructure:"deposit"`
	// Withdraw defines the withdraw event hash
	Withdraw string `mapstructure:"withdraw"`
}

const (
	BitcoinRPCConfigFileName = "bitcoin.toml"
)

func LoadBitcoinConfig(homePath string) (*BitconConfig, error) {
	config := BitconConfig{}
	configFile := path.Join(homePath, BitcoinRPCConfigFileName)
	_, err := os.Stat(configFile)
	if err != nil {
		// if file not exist use default config
		// TODO: add gen config command after, The default configuration may not be required
		if !os.IsNotExist(err) {
			return nil, err
		}
		// config.NetworkName = "signet"
		// config.RPCHost = "localhost"
		// config.RPCPort = "8332"
		// config.RPCUser = "user"
		// config.RPCPass = "password"
		// config.WalletName = "walletname"
		// config.Destination = "tb1qgm39cu009lyvq93afx47pp4h9wxq5x92lxxgnz"
		// config.IndexerListenAddress = "tb1qsja4hvx66jr9grgmt8452letmz37gmludcrnup"
		// config.EnableIndexer = false
		// config.EnableCommitter = false
		// return &config, nil
	}

	v := viper.New()
	v.SetConfigFile(configFile)
	v.SetEnvPrefix("BITCOIN")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		// Load from environment variables if not found in the config file
		err = v.BindEnv("network-name", "BITCOIN_NETWORK_NAME")
		if err != nil {
			return nil, err
		}
		err = v.BindEnv("rpc-host", "BITCOIN_RPC_HOST")
		if err != nil {
			return nil, err
		}
		err = v.BindEnv("rpc-port", "BITCOIN_RPC_PORT")
		if err != nil {
			return nil, err
		}
		err = v.BindEnv("rpc-user", "BITCOIN_RPC_USER")
		if err != nil {
			return nil, err
		}
		err = v.BindEnv("rpc-pass", "BITCOIN_RPC_PASS")
		if err != nil {
			return nil, err
		}
		err = v.BindEnv("wallet-name", "BITCOIN_WALLET_NAME")
		if err != nil {
			return nil, err
		}
		err = v.BindEnv("destination", "BITCOIN_DESTINATION")
		if err != nil {
			return nil, err
		}
		err = v.BindEnv("enable-indexer", "BITCOIN_ENABLE_INDEXER")
		if err != nil {
			return nil, err
		}
		err = v.BindEnv("enable-committer", "BITCOIN_ENABLE_COMMITTER")
		if err != nil {
			return nil, err
		}
		err = v.BindEnv("indexer-listen-address", "BITCOIN_INDEXER_LISTEN_ADDRESS")
		if err != nil {
			return nil, err
		}
		err = v.BindEnv("bridge.eth-rpc-url", "BITCOIN_BRIDGE_ETH_RPC_URL")
		if err != nil {
			return nil, err
		}
		err = v.BindEnv("bridge.eth-priv-key", "BITCOIN_BRIDGE_ETH_PRIV_KEY")
		if err != nil {
			return nil, err
		}
		err = v.BindEnv("bridge.contract-address", "BITCOIN_BRIDGE_CONTRACT_ADDRESS")
		if err != nil {
			return nil, err
		}
		err = v.BindEnv("bridge.abi", "BITCOIN_BRIDGE_ABI")
		if err != nil {
			return nil, err
		}
		err = v.BindEnv("bridge.gas-limit", "BITCOIN_BRIDGE_GAS_LIMIT")
		if err != nil {
			return nil, err
		}
		err = v.BindEnv("bridge.aa-sca-registry", "BITCOIN_BRIDGE_AA_SCA_REGISTRY")
		if err != nil {
			return nil, err
		}
		err = v.BindEnv("bridge.aa-kernel-factory", "BITCOIN_BRIDGE_AA_KERNEL_FACTORY")
		if err != nil {
			return nil, err
		}
		err = v.BindEnv("state.host", "BITCOIN_STATE_HOST")
		if err != nil {
			return nil, err
		}
		err = v.BindEnv("state.port", "BITCOIN_STATE_PORT")
		if err != nil {
			return nil, err
		}
		err = v.BindEnv("state.user", "BITCOIN_STATE_USER")
		if err != nil {
			return nil, err
		}
		err = v.BindEnv("state.pass", "BITCOIN_STATE_PASS")
		if err != nil {
			return nil, err
		}
		err = v.BindEnv("state.db-name", "BITCOIN_STATE_DB_NAME")
		if err != nil {
			return nil, err
		}
		err = v.BindEnv("fee", "BITCOIN_FEE")
		if err != nil {
			return nil, err
		}
		err = v.BindEnv("evm.enable-listener", "BITCOIN_EVM_ENABLE_LISTENER")
		if err != nil {
			return nil, err
		}
		err = v.BindEnv("evm.deposit", "BITCOIN_EVM_DEPOSIT")
		if err != nil {
			return nil, err
		}
		err = v.BindEnv("evm.withdraw", "BITCOIN_EVM_WITHDRAW")
		if err != nil {
			return nil, err
		}
	}

	err = v.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// ChainParams get chain params by network name
func ChainParams(network string) *chaincfg.Params {
	switch network {
	case MAINNET:
		return &chaincfg.MainNetParams
	case TESTNET:
		return &chaincfg.TestNet3Params
	case SIGNET:
		return &chaincfg.SigNetParams
	case SIMNET:
		return &chaincfg.SimNetParams
	case REGTEST:
		return &chaincfg.RegressionNetParams
	default:
		return &chaincfg.TestNet3Params
	}
}
