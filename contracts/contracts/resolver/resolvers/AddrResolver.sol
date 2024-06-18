// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./IAddrResolver.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract AddrResolver is IAddrResolver {
    address private _addr;

    constructor(address resolveToAddr) {
        _addr = resolveToAddr;
    }

    /**
     * @dev Returns the address that this resolver resolves to.
     */
    function addr() public view returns (address) {
        return _addr;
    }
}
