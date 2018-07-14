# Geth Server Prometheus Exporter
Monitor your Ethereum Geth server with Prometheus and Grafana. Checkout the [Grafana Dashboard](https://grafana.com/dashboards/6976) to implement a beautiful geth server monitor for your own server, or you can just import Dashboard ID: `6976` once you have GethExporter up and running.

<p align="center"><img width="90%" src="https://img.cjx.io/gethexporter-grafana.png"></p>

## Docker
Run this Prometheus Exporter in a [Docker container](https://hub.docker.com/r/hunterlong/gethexporter/builds/)! Include your Geth server endpoint as `GETH` environment variable. 
```bash
docker run -it -d -p 9090:9090 \
  -e "GETH="http://mygethserverhere.com:8545" \
  hunterlong/gethexporter
```

## Features
- Current and Average Gas Price
- Total amount of ERC20 Token Transfers
- Total amount of ETH transactions
- Pending Transaction count

## Prometheus Response
```
geth_block 5959116
geth_seconds_last_block 140.34
geth_block_transactions 33
geth_block_value 9.531768449804128
geth_block_gas_used 7983611
geth_block_gas_limit 7999992
geth_block_nonce 15640153714112430600
geth_block_difficulty 3333727730300686
geth_block_uncles 1
geth_block_size_bytes 20450
geth_gas_price 41000000000
geth_pending_transactions 186
geth_network_id 1
geth_contracts_created 0
geth_token_transfers 6
geth_eth_transfers 16
```
