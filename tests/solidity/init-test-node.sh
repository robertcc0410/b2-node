set -x
set -e

CHAINID="bsquarednetwork_102-11"
MONIKER="b2network"

# localKey address 0xf4306d09e40dD74A6B849c69007D42a7f4806114
GENESIS_KEY="genesis"

# TODO: populate this with the mnemonic for the genesis account
# GENESIS_MNEMONIC=""
GENESIS_BALANCE="1000000000000000000000stake"
GLOBAL_BALANCE="1000000000000000000000aphoton"

importKey() {
    # Import keys from mnemonics
    echo $GENESIS_MNEMONIC | ethermintd keys add $GENESIS_KEY --recover --keyring-backend test --algo "eth_secp256k1"
}

updateConf() {
    # Set gas limit in genesis
    cat $HOME/.ethermintd/config/genesis.json | jq '.consensus_params["block"]["max_gas"]="10000000"' >$HOME/.ethermintd/config/tmp_genesis.json
    mv $HOME/.ethermintd/config/tmp_genesis.json $HOME/.ethermintd/config/genesis.json

    find $HOME/.ethermintd/config -name 'config.toml' -exec toml set --toml-path {} --to-bool consensus.create_empty_blocks true \;
    find $HOME/.ethermintd/config -name 'config.toml' -exec toml set --toml-path {} --to-bool consensus.create_empty_blocks true \;
    find $HOME/.ethermintd/config -name 'config.toml' -exec toml set --toml-path {} consensus.timeout_commit 10s \;
    find $HOME/.ethermintd/config -name 'config.toml' -exec toml set --toml-path {} rpc.laddr tcp://0.0.0.0:26657 \;
    find $HOME/.ethermintd/config -name 'config.toml' -exec toml set --toml-path {} --to-bool instrumentation.prometheus true \;

    find $HOME/.ethermintd/config -name 'app.toml' -exec toml set --toml-path {} --to-int telemetry.prometheus-retention-time 1000000000000 \;
    find $HOME/.ethermintd/config -name 'app.toml' -exec toml set --toml-path {} --to-bool api.enabled true \;
    find $HOME/.ethermintd/config -name 'app.toml' -exec toml set --toml-path {} json-rpc.api eth,txpool,net,web3 \;
    find $HOME/.ethermintd/config -name 'app.toml' -exec toml set --toml-path {} json-rpc.address 0.0.0.0:8545 \;
    find $HOME/.ethermintd/config -name 'app.toml' -exec toml set --toml-path {} json-rpc.ws-address 0.0.0.0:8546 \;
    find $HOME/.ethermintd/config -name 'app.toml' -exec toml set --toml-path {} json-rpc.metrics-address 0.0.0.0:6065 \;
}

allocateAccount() {
    # Allocate genesis accounts (cosmos formatted addresses)
    ethermintd add-genesis-account "$(ethermintd keys show $GENESIS_KEY -a --keyring-backend test)" "$GLOBAL_BALANCE,$GENESIS_BALANCE" --keyring-backend test
}

init() {
    # remove existing daemon and client
    rm -rf \
        /root/.ethermintd/config \
        /root/.ethermintd/data \
        /root/.ethermintd/keyring-test

    importKey
    # git init $HOME/.ethermintd
    git -C $HOME/.ethermintd add .
    git -C $HOME/.ethermintd commit --message "import key"

    ethermintd init $MONIKER \
        --chain-id $CHAINID
    git -C $HOME/.ethermintd add .
    git -C $HOME/.ethermintd commit --message "init chain"

    updateConf
    git -C $HOME/.ethermintd add .
    git -C $HOME/.ethermintd commit --message "update config"

    allocateAccount
    git -C $HOME/.ethermintd add .
    git -C $HOME/.ethermintd commit --message "allocate account"

    ethermintd gentx $GENESIS_KEY $GENESIS_BALANCE \
        --amount=$GLOBAL_BALANCE \
        --chain-id $CHAINID \
        --keyring-backend test
    git -C $HOME/.ethermintd add .
    git -C $HOME/.ethermintd commit --message "sign genesis transaction"

    ethermintd collect-gentxs
    git -C $HOME/.ethermintd add .
    git -C $HOME/.ethermintd commit --message "collect genesis tx"

    set -e
    ethermintd validate-genesis
    git -C $HOME/.ethermintd add .
    # git -C $HOME/.ethermintd commit --message "run this to ensure everything worked and that the genesis file is setup correctly"

    # Start the node (remove the --pruning=nothing flag if historical queries are not needed)
    # ethermintd start --metrics --pruning=nothing --rpc.unsafe --keyring-backend test --log_level info --json-rpc.api eth,txpool,personal,net,debug,web3 --api.enable
}

$@
