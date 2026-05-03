// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {VaultAccess} from "./VaultAccess.sol";
import {IVaultFunds} from "../interfaces/IVaultFunds.sol";

abstract contract VaultFunds is VaultAccess, IVaultFunds {
    function save() external payable onlyAccount {
        _save(msg.sender, msg.value);
    }

    receive() external payable {
        _save(msg.sender, msg.value);
    }

    function withdraw(uint256 amount) external onlyAccount {
        if (amount == 0) {
            revert InvalidAmount();
        }

        uint256 available = saved[msg.sender];

        if (available < amount) {
            revert InsufficientSavedBalance(msg.sender, amount, available);
        }

        saved[msg.sender] = available - amount;
        (bool success, ) = payable(msg.sender).call{value: amount}("");

        if (!success) {
            revert WithdrawFailed(msg.sender, amount);
        }

        emit Withdrawn(msg.sender, amount);
    }

    function withdrawAll() external onlyAccount {
        uint256 amount = saved[msg.sender];

        if (amount == 0) {
            revert NothingToWithdraw(msg.sender);
        }

        saved[msg.sender] = 0;
        (bool success, ) = payable(msg.sender).call{value: amount}("");

        if (!success) {
            revert WithdrawFailed(msg.sender, amount);
        }

        emit Withdrawn(msg.sender, amount);
    }

    function getVaultBalance() external view returns (uint256) {
        return address(this).balance;
    }

    function _save(address account, uint256 amount) internal {
        if (!isAccount[account]) {
            revert NotRegisteredAccount(account);
        }

        if (amount == 0) {
            revert NoEthSent();
        }

        saved[account] += amount;

        emit Saved(account, amount);
    }
}
