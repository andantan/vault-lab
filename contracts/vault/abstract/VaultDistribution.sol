// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {VaultAccess} from "./VaultAccess.sol";
import {IVaultDistribution} from "../interfaces/IVaultDistribution.sol";

abstract contract VaultDistribution is VaultAccess, IVaultDistribution {
    function spread() external onlyAccount {
        uint256 totalAmount = saved[msg.sender];

        if (totalAmount == 0) {
            revert NothingToSpread(msg.sender);
        }

        uint256 recipientCount = accountList.length - 1;

        if (recipientCount == 0) {
            revert AmountTooSmall(totalAmount, recipientCount);
        }

        uint256 share = totalAmount / recipientCount;
        uint256 remainder = totalAmount % recipientCount;

        if (share == 0) {
            revert AmountTooSmall(totalAmount, recipientCount);
        }

        saved[msg.sender] = remainder;

        for (uint256 i = 0; i < accountList.length; i++) {
            address account = accountList[i];

            if (account != msg.sender) {
                saved[account] += share;
            }
        }

        emit Spread(msg.sender, totalAmount, share, remainder);
    }

    function collect() external onlyMaster {
        uint256 totalAmount = 0;

        for (uint256 i = 0; i < accountList.length; i++) {
            address account = accountList[i];

            if (account != master) {
                totalAmount += saved[account];
                saved[account] = 0;
            }
        }

        if (totalAmount == 0) {
            revert NothingToCollect();
        }

        saved[master] += totalAmount;

        emit Collected(master, totalAmount);
    }
}