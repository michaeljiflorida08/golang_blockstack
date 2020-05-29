# golang_blockstack

#launch test net and output to log file, name it as "trace_logging.txt"

RUST_BACKTRACE=FULL BLOCKSTACK_DEBUG=1 cargo testnet start --config=./testnet/stacks-node/conf/neon-miner-conf.toml |& tee -a trace_logging.txt

#run golang command line 

go run golang_blockstack.go 2
