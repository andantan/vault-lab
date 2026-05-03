#!/usr/bin/env bash
set -euo pipefail

OUT="build/standard-json-input.json"

mkdir -p build

node <<'NODE'
const fs = require("fs");
const path = require("path");

const sources = {};
const files = [
  "contracts/MultiAccountVault.sol",
  "contracts/abstract/VaultAccess.sol",
  "contracts/abstract/VaultAccounts.sol",
  "contracts/abstract/VaultFunds.sol",
  "contracts/abstract/VaultDistribution.sol",
  "contracts/interfaces/IMultiAccountVault.sol",
  "contracts/interfaces/IVaultAccounts.sol",
  "contracts/interfaces/IVaultFunds.sol",
  "contracts/interfaces/IVaultDistribution.sol",
  "contracts/libraries/VaultTypes.sol"
];

for (const file of files) {
  sources[file] = {
    content: fs.readFileSync(file, "utf8")
  };
}

const input = {
  language: "Solidity",
  sources,
  settings: {
    optimizer: {
      enabled: false,
      runs: 200
    },
    outputSelection: {
      "*": {
        "*": [
          "abi",
          "evm.bytecode",
          "evm.deployedBytecode",
          "metadata"
        ]
      }
    }
  }
};

fs.writeFileSync("build/standard-json-input.json", JSON.stringify(input, null, 2));
NODE

echo "wrote ${OUT}"