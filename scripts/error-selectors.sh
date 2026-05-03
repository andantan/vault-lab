#!/usr/bin/env bash
set -euo pipefail

TARGET="${1:-}"

node - "$TARGET" <<'NODE'
const { id } = require("ethers");

const target = process.argv[2];

const errors = [
  "OnlyMaster(address,address)",
  "NotRegisteredAccount(address)",
  "AlreadyRegisteredAccount(address)",
  "InvalidAccount()",
  "InvalidRole(uint8)",
  "CannotRemoveMaster()",
  "AccountHasSavedBalance(address,uint256)",
  "AccountIndexOutOfBounds(uint256,uint256)",
  "NoEthSent()",
  "InvalidAmount()",
  "InsufficientSavedBalance(address,uint256,uint256)",
  "NothingToWithdraw(address)",
  "WithdrawFailed(address,uint256)",
  "NothingToSpread(address)",
  "AmountTooSmall(uint256,uint256)",
  "NothingToCollect()"
];

let matched = false;

for (const signature of errors) {
  const selector = id(signature).slice(0, 10);

  if (!target || selector === target) {
    console.log(`${selector}  ${signature}`);
    matched = true;
  }
}

if (target && !matched) {
  console.error(`No matching error selector found: ${target}`);
  process.exit(1);
}
NODE