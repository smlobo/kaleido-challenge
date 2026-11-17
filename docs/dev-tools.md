# Developer tools

Commands for inspecting the local Besu chain and working with the contracts.
All examples assume the dev stack from `docker/docker-compose.yaml` (Besu and Postgres) is running:

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

Note: in Foundry 1.9, `cast block --json`, `cast tx --json` and `cast wallet new --json`
wrap the result in `{"data": ...}`. `cast receipt --json` doesn't.

### Transactions

```bash
TX=0xbebc550c644ccd42991fd6c7f8defdbcfa56fd81927723c23f392c29af714683

cast tx $TX                      # full transaction
cast tx $TX type                 # 0x2 = EIP-1559
cast tx $TX nonce
cast tx $TX blockNumber
```

Decode what a transaction called:

```bash
cast 4byte-decode $(cast tx $TX input)                               # looks up the selector online
cast calldata-decode "safeMint(address,uint256)" $(cast tx $TX input)  # offline, with a known signature
```

### Receipts

The transaction (`cast tx`) is what was signed and sent. The receipt is what
happened when it ran: success or failure, gas used, the new contract address, and the
events it emitted. QBFT blocks are final as soon as they're made, so one receipt is
the final answer. This is what `chain.WaitMined` polls for.

```bash
DEPLOY=0xec114cc9aab4d3f1f85ec1fd728ddc10340c272af3da7426a8d5c6d83215ed2d
MINT=0x58afc1c2d15d59513393bacf1f443e88f30304cb25b6fc20917375ed0aa10c0a

cast receipt $DEPLOY                  # all fields
cast receipt $DEPLOY status           # true / false
cast receipt $DEPLOY contractAddress  # set only for a contract creation
cast receipt $DEPLOY gasUsed          # actual gas, at most the tx's gasLimit
cast receipt $MINT logs               # events emitted, as JSON
cast receipt $MINT --json | jq '.logs[0].topics'
```

`cast receipt` waits until the transaction is mined. Add `--async` to return
immediately for a transaction that's still pending.

Fields for the `m1deploy` deploy above:

| Field | Example | Meaning |
|---|---|---|
| `status` | `1 (success)` | `0` means it reverted. All state changes are undone, but it still uses a nonce and appears in the block |
| `contractAddress` | `0xfb2c…f4d7` | Address of the new contract. Empty for a normal call |
| `gasUsed` | 1,052,698 | Gas this transaction actually used. `m1deploy` sets `gasLimit` to the estimate + 20% (1,274,664) |
| `cumulativeGasUsed` | 1,052,698 | Gas used by this and all earlier transactions in the block. Equal to `gasUsed` when it's first in the block |
| `effectiveGasPrice` | 0 | Price paid per unit of gas. Always 0 on this chain |
| `logs` | 1 entry | Events. The deploy emits `OwnershipTransferred(0x0 → deployer)` from the `Ownable` constructor |
| `logsBloom` | `0x0000…` | A compact filter that lets `eth_getLogs` skip blocks that can't contain a matching event |
| `root`, `blobGasUsed`, `blobGasPrice` | empty | Not used here (pre-Byzantium state root; EIP-4844 blob transactions) |

#### Reading logs

Each log has the emitting contract's `address`, up to 4 `topics` and `data`. `topics[0]` is
the event signature hash. The rest are the `indexed` parameters, each padded to 32 bytes.
The mint's `Transfer` log:

```
address  0xfb2c4d7720fd5ec5c579ef7dbc7ef08d92a4f4d7   the DemoNFT
topics   0xddf252ad…b3ef   keccak("Transfer(address,address,uint256)")
         0x0000…0000       from    = zero address, so it's a mint
         0x0000…7223e0a    to      = 0x05cf9cedfd7f42eae4b768b824c01efaf7223e0a
         0x0000…0001       tokenId = 1
data     0x                nothing un-indexed
```

To get the signature hash of any event:

```bash
cast keccak "Transfer(address,address,uint256)"
```

#### Failed transactions

Besu runs with `--revert-reason-enabled`, so a failed receipt carries the encoded
error. Example: a stranger calling `safeMint` with an explicit gas limit, so that gas
estimation doesn't catch the revert before sending:

```bash
FAILED=0x0dcf127cfbfa2d68056de3d6c6e115ed3e2fa16ade30b89c16897d76b8937612

cast receipt $FAILED --json | jq '{status, revertReason, gasUsed}'
# {"status":"0x0","revertReason":"0x118cdaa7000…755a9774…","gasUsed":"0x5db6"}

cast 4byte 0x118cdaa7         # OwnableUnauthorizedAccount(address)
cast decode-error 0x118cdaa7000000000000000000000000755a9774ec3f90ea08abd61076ccd63bfacba6d6 \
  --sig "OwnableUnauthorizedAccount(address)"

cast run $FAILED              # replays it with a call trace ending in the revert
```

The first 4 bytes of `revertReason` are the custom error's selector, and the rest are its
ABI-encoded arguments. The Go code rarely sees a failed receipt, because
`eth_estimateGas` simulates the call first and returns the revert before anything is
signed. A transaction can still fail on-chain if the state changes between the
estimate and the block, for example if another buyer takes the last ticket. That's the
reason `WaitMined` checks `status`.

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

## Postgres

The compose Postgres listens on host port **5432**. It has no volume: `docker compose down` 
wipes it together with the chain, so stored keys and nonces never outlive the chain they belong to.

```bash
export PGPASSWORD=ticketing
psql -h localhost -p 5432 -U ticketing -d ticketing
```

Useful queries:

```sql
\dt                                         -- tables
SELECT version FROM schema_migrations;      -- applied migrations
SELECT purpose, address, next_nonce, length(encrypted_private_key) AS sealed_bytes FROM signing_keys; 
-- 61 bytes = version + GCM nonce + key + tag
SELECT nonce, status, block_number, contract_address, tx_hash FROM transactions ORDER BY created_at;
-- the outbox
```

Check that a key's database nonce matches the chain:

```bash
cast nonce <address>    # should equal signing_keys.next_nonce once everything is mined
```

If the database is ahead, the chain was reset without wiping Postgres. Reset both:

```bash
docker compose -f docker/docker-compose.yaml down
```

```bash
docker compose -f docker/docker-compose.yaml up -d --wait
```

Tests don't touch the `public` schema. Each one creates its own `test_<random>`
schema, migrates it and drops it afterwards. Point them at another database with
`DATABASE_URL`.

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
KEY_ENCRYPTION_KEY=... go run ./cmd/m2demo   # same, with the key and nonces in Postgres
go test ./...            # Besu / Postgres tests skip if those aren't running
BESU_RPC_URL=http://other-host:8545 go test ./internal/chain/
```
