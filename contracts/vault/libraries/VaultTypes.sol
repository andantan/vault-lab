// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

library VaultTypes {
    enum AccountRole {
        Unknown,
        Master,
        Slave
    }

    struct AccountInfo {
        address account;
        AccountRole role;
        uint256 saved;
        bool registered;
    }

    function isManagedRole(AccountRole role) internal pure returns (bool) {
        return role == AccountRole.Master || role == AccountRole.Slave;
    }

    function isMaster(AccountRole role) internal pure returns (bool) {
        return role == AccountRole.Master;
    }

    function isParticipant(AccountRole role) internal pure returns (bool) {
        return role == AccountRole.Master || role == AccountRole.Slave;
    }
}
