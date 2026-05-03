// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {IVaultAccounts} from "./IVaultAccounts.sol";
import {IVaultFunds} from "./IVaultFunds.sol";
import {IVaultDistribution} from "./IVaultDistribution.sol";

interface IMultiAccountVault is
    IVaultAccounts,
    IVaultFunds,
    IVaultDistribution
{}