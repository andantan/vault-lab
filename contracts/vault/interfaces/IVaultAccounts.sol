// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {VaultTypes} from "../libraries/VaultTypes.sol";

interface IVaultAccounts {
    error AlreadyRegisteredAccount(address account);
    error InvalidAccount();
    error InvalidRole(VaultTypes.AccountRole role);
    error CannotRemoveMaster();
    error AccountHasSavedBalance(address account, uint256 savedAmount);
    error AccountIndexOutOfBounds(uint256 index, uint256 length);

    event AccountAdded(address indexed account, VaultTypes.AccountRole role);
    event AccountRemoved(address indexed account);

    function addAccount(address account, VaultTypes.AccountRole role) external;
    function removeAccount(address account) external;
    function getAccounts() external view returns (address[] memory);
    function accountAt(uint256 index) external view returns (address);
    function accountCount() external view returns (uint256);
    function getAccountInfo(address account) external view returns (VaultTypes.AccountInfo memory);
}