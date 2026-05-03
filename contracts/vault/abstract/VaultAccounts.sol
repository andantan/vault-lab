// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {VaultAccess} from "./VaultAccess.sol";
import {IVaultAccounts} from "../interfaces/IVaultAccounts.sol";
import {VaultTypes} from "../libraries/VaultTypes.sol";

abstract contract VaultAccounts is VaultAccess, IVaultAccounts {
    function addAccount(address account, VaultTypes.AccountRole role) external onlyMaster {
        _addAccount(account, role);
    }

    function removeAccount(address account) external onlyMaster {
        _removeAccount(account);
    }

    // Warning: this copies the entire account list to memory
    // and can be expensive for large lists.
    function getAccounts() external view returns (address[] memory) {
        return accountList;
    }

    function accountAt(uint256 index) external view returns (address) {
        uint256 length = accountList.length;

        if (index >= length) {
            revert AccountIndexOutOfBounds(index, length);
        }

        return accountList[index];
    }

    function accountCount() external view returns (uint256) {
        return accountList.length;
    }

    function getAccountInfo(address account) external view returns (VaultTypes.AccountInfo memory) {
        return VaultTypes.AccountInfo({
            account: account,
            role: roles[account],
            saved: saved[account],
            registered: isAccount[account]
        });
    }

    function _addAccount(address account, VaultTypes.AccountRole role) internal {
        if (account == address(0)) {
            revert InvalidAccount();
        }

        if (isAccount[account]) {
            revert AlreadyRegisteredAccount(account);
        }

        if (!VaultTypes.isManagedRole(role)) {
            revert InvalidRole(role);
        }

        if (role == VaultTypes.AccountRole.Master) {
            if (master != address(0)) {
                revert InvalidRole(role);
            }

            master = account;
        }

        isAccount[account] = true;
        roles[account] = role;
        accountIndex[account] = accountList.length;
        accountList.push(account);

        emit AccountAdded(account, role);
    }

    function _removeAccount(address account) internal {
        if (account == master) {
            revert CannotRemoveMaster();
        }

        if (!isAccount[account]) {
            revert NotRegisteredAccount(account);
        }

        if (saved[account] != 0) {
            revert AccountHasSavedBalance(account, saved[account]);
        }

        uint256 index = accountIndex[account];
        uint256 lastIndex = accountList.length - 1;

        if (index != lastIndex) {
            address lastAccount = accountList[lastIndex];
            accountList[index] = lastAccount;
            accountIndex[lastAccount] = index;
        }

        accountList.pop();

        delete accountIndex[account];
        delete isAccount[account];
        delete roles[account];
        delete saved[account];

        emit AccountRemoved(account);
    }
}
