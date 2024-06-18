// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./IServerResolver.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract ServerResolver is IServerResolver {
    string private _endpoint;
    
    constructor(string memory endpoint) {
        _endpoint = endpoint;
    }

    /**
     * @dev Returns the server endpoint that this resolver resolves to.
     */
    function server() public view returns (string memory) {
        return _endpoint;
    }
}
