// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

interface IVaultFunds {
    error NoEthSent();
    error InvalidAmount();
    error InsufficientSavedBalance(address account, uint256 requested, uint256 available);
    error NothingToWithdraw(address account);
    error WithdrawFailed(address account, uint256 amount);

    event Saved(address indexed account, uint256 amount);
    event Withdrawn(address indexed account, uint256 amount);

    function save() external payable;
    function withdraw(uint256 amount) external;
    function withdrawAll() external;
    function getVaultBalance() external view returns (uint256);
}