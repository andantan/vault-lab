// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {VaultAccess} from "./abstract/VaultAccess.sol";
import {VaultAccounts} from "./abstract/VaultAccounts.sol";
import {VaultDistribution} from "./abstract/VaultDistribution.sol";
import {VaultFunds} from "./abstract/VaultFunds.sol";
import {VaultTypes} from "./libraries/VaultTypes.sol";

contract MultiAccountVault is
    VaultAccounts,
    VaultFunds,
    VaultDistribution
{
    constructor() {
        _addAccount(msg.sender, VaultTypes.AccountRole.Master);
    }
}