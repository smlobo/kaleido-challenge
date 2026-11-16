#!/usr/bin/env bash
# Compiles contracts with Foundry (forge build; solc 0.8.38 per foundry.toml)
# and regenerates the Go bindings. The ABI/bytecode in build/ and the bindings
# in internal/contracts are committed, so building the service needs neither
# Foundry nor Node.
set -euo pipefail
cd "$(dirname "$0")"

command -v forge >/dev/null || { echo "forge not found: brew install foundry" >&2; exit 1; }
command -v jq >/dev/null || { echo "jq not found: brew install jq" >&2; exit 1; }
[ -d node_modules/@openzeppelin ] || npm install --silent

ABIGEN="go run github.com/ethereum/go-ethereum/cmd/abigen@v1.17.10"

forge build

export_contract() { # export_contract <ContractName> <go-package>
  local name=$1 pkg=$2 artifact="out/$1.sol/$1.json"
  mkdir -p build "../internal/contracts/$pkg"
  jq -c '.abi' "$artifact" > "build/$name.abi"
  jq -r '.bytecode.object' "$artifact" | sed 's/^0x//' > "build/$name.bin"
  (cd .. && $ABIGEN --v2 --abi "contracts/build/$name.abi" --bin "contracts/build/$name.bin" \
    --pkg "$pkg" --type "$name" --out "internal/contracts/$pkg/$pkg.go")
  echo "built $name -> internal/contracts/$pkg"
}

export_contract DemoNFT demonft
