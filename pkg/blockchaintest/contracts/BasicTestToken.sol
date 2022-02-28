// contracts/BasicTestToken.sol
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.11;

import "./ERC20.sol";

contract BasicTestToken is ERC20 {
    constructor() ERC20("Basic test token", "BTTX") {
        _mint(msg.sender, uint256(20000000000000000000000000)); // 20 mil tokens
    }
}