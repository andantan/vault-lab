// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {VaultTypes} from "../libraries/VaultTypes.sol";

abstract contract VaultAccess {
    address public master;
    address[] internal accountList;

    mapping(address => bool) public isAccount;
    mapping(address => VaultTypes.AccountRole) public roles;
    mapping(address => uint256) public saved;
    mapping(address => uint256) internal accountIndex;

    error OnlyMaster(address caller, address master);
    error NotRegisteredAccount(address account);

    modifier onlyMaster() {
        if (msg.sender != master) {
            revert OnlyMaster(msg.sender, master);
        }

        _;
    }

    modifier onlyAccount() {
        if (!isAccount[msg.sender]) {
            revert NotRegisteredAccount(msg.sender);
        }

        _;
    }
}