# Developer tools

Commands for inspecting the local Besu chain and working with the contracts.
All examples assume the dev chain from `docker/docker-compose.yaml` is running:

```bash
docker compose -f docker/docker-compose.yaml up -d --wait
```

## cast (Foundry)

Install Foundry with `brew install foundry`. Set the RPC URL once per shell
so you can leave out `--rpc-url` on every command:

```bash
export ETH_RPC_URL=http://localhost:8545
```

The examples use addresses and hashes from a real `m1deploy` run. Swap in the
ones from your own run. The chain resets on `docker compose down`.

### Chain

```bash
cast chain-id                 # 1337
cast block-number             # latest block height (blocks = height + 1, counting genesis)
cast block 0 -f hash          # genesis hash; it changes when the chain is reset
cast block 1039 -f timestamp  # one field of a block
cast block 1039 --json | jq '.data.transactions'   # tx hashes in a block
```

Note: `cast --json` output in Foundry 1.9 wraps the result in `{"data": ...}`.

### Transactions and receipts

```bash
TX=0xbebc550c644ccd42991fd6c7f8defdbcfa56fd81927723c23f392c29af714683

cast tx $TX                      # full transaction
cast tx $TX type                 # 0x2 = EIP-1559
cast tx $TX nonce
cast tx $TX blockNumber
cast receipt $TX status          # true = success
cast receipt 0xe83fe5a0734ec78f5bb5a7f5aa8e76afe7940892901271a3d56b6eb013ac6280 contractAddress
```

Decode what a transaction called:

```bash
cast 4byte-decode $(cast tx $TX input)                               # looks up the selector online
cast calldata-decode "safeMint(address,uint256)" $(cast tx $TX input)  # offline, with a known signature
```

### Reading contracts

`cast call` runs an `eth_call`. It needs no key and sends no transaction. Write the
return type in the second set of parentheses to get a decoded result.

```bash
NFT=0x9c4faec1a6dc6331d1408b15d7671edccb35d9c5

cast call $NFT "name()(string)"
cast call $NFT "symbol()(string)"
cast call $NFT "ownerOf(uint256)(address)" 1
cast call $NFT "balanceOf(address)(uint256)" 0xbf9c35b275b3b5406fb7b23978af3990cb69f5b9
cast codesize $NFT               # 0 means no contract at that address
```

### Accounts

```bash
cast balance 0xc06d17055c0a9eee97d0e08d084a3bc39572230d   # 0: gas is free, keys need no funding
cast nonce   0xc06d17055c0a9eee97d0e08d084a3bc39572230d   # transactions sent so far
```

### Raw JSON-RPC

For any method without a dedicated command:

```bash
cast rpc eth_getBlockByNumber latest false
cast rpc txpool_besuStatistics
```

## Without cast (curl)

```bash
curl -s -X POST -H 'content-type: application/json' \
  --data '{"jsonrpc":"2.0","id":1,"method":"eth_blockNumber","params":[]}' \
  localhost:8545
```

The result is hex, for example `"0x4aa"` = 1194.

## Besu logs

The node logs one `Produced #N` line per block, about every 2 seconds, with the
number of transactions it included:

```bash
docker logs docker-besu-1 2>&1 | grep Produced | tail                # recent blocks
docker logs docker-besu-1 2>&1 | grep Produced | grep -v '|    0 tx' # only blocks with transactions
docker logs -f docker-besu-1 2>&1 | grep Produced                    # follow live
```

The logs don't show what a transaction did. Use `cast tx` / `cast receipt`
or `scripts/find_deploys.py` for that.

`Unable to find sync target. Waiting for 5 peers` is harmless. It's a
single-node network with discovery off.

## scripts/find_deploys.py

Lists every deployment of a compiled contract. It scans each block for creation
transactions (no `to`) whose input starts with the contract's bytecode from
`contracts/build/<Name>.bin`, then prints the new address and the deployed
contract's `name()` / `symbol()`.

```bash
python3 scripts/find_deploys.py                  # DemoNFT
python3 scripts/find_deploys.py TicketNFT        # another contract in contracts/build
python3 scripts/find_deploys.py --rpc http://localhost:8545 --batch 500
```

Output:

```
block 2205  tx 0xed5e…  from 0x2971…  contract 0x6755…  ok  'Demo Ticket' (DEMO)
```

Notes:
- `m1deploy` deploys are named `Demo Ticket`. `go test` deploys are named `Test`,
  two per run.
- It matches the bytecode *as currently built*. After a contract is rebuilt with
  different settings, or recompiled in a way that changes its metadata hash, older
  deploys no longer match.
- It uses only the Python standard library. Blocks are fetched in batched JSON-RPC requests.

## scripts/find_minted.py

Lists every ERC-721 mint on the chain, from any contract, with each token's current
owner. A mint is a `Transfer` event from the zero address, so this uses
`eth_getLogs` instead of reading every block. It takes well under a second.

```bash
python3 scripts/find_minted.py                                      # all mints
python3 scripts/find_minted.py --contract 0x9c4faec1a6dc6331d1408b15d7671edccb35d9c5
python3 scripts/find_minted.py --to 0x000000000000000000000000000000000000bEEF
python3 scripts/find_minted.py --from-block 2000
```

Output:

```
block 1039  tx 0xbebc…  contract 0x9c4f… 'Demo Ticket' (DEMO)  token 1  to 0xbf9c…
```

If the token has since moved, the line ends with `now <owner>`, or `now burned`
if `ownerOf` reverts. That will show resales and check-ins once `TicketNFT` exists.

Notes:
- ERC-20 `Transfer` events have the same signature hash. The script skips them by
  keeping only logs with 4 topics, because ERC-721 indexes `tokenId` and ERC-20 doesn't.
- It doesn't depend on bytecode, so unlike `find_deploys.py` it still finds mints on
  contracts built before a rebuild.
- Besu caps one `eth_getLogs` call at 5000 blocks (`--rpc-max-logs-range`), so the
  script queries in 5000-block chunks.
- The same query with cast (the topic is `cast keccak "Transfer(address,address,uint256)"`):

  ```bash
  cast logs --from-block 0 "Transfer(address indexed,address indexed,uint256 indexed)" \
    0x0000000000000000000000000000000000000000
  ```

Both scripts share their JSON-RPC helpers through `scripts/ethrpc.py`.

## contracts/build.sh

Compiles the contracts with Forge and regenerates the Go bindings. Run it from
anywhere:

```bash
./contracts/build.sh
```

It:
1. runs `npm install` if OpenZeppelin is missing from `contracts/node_modules`
2. runs `forge build` (solc 0.8.38, settings in `contracts/foundry.toml`)
3. copies the ABI and bytecode from `contracts/out/<Name>.sol/<Name>.json` into
   `contracts/build/<Name>.abi` / `.bin` (committed)
4. runs abigen v2 to write `internal/contracts/<pkg>/<pkg>.go` (committed)

Requires `forge` and `jq`. You only need it after changing a contract. Building
the Go service uses the committed bindings.

## Go smoke test and tests

```bash
go run ./cmd/m1deploy    # deploy DemoNFT, mint token 1, verify ownerOf
go test ./...            # the Besu tests skip if the chain is not running
BESU_RPC_URL=http://other-host:8545 go test ./internal/chain/
```
