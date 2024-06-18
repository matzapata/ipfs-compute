// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./IRsaPublicKeyResolver.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract RsaPublicKeyResolver is IRsaPublicKeyResolver {
    string private _rsaPublicKey;

    constructor(string memory publicKey)  {
        _rsaPublicKey = publicKey;
    }

    /**
     * @dev Returns the RSA public key that this resolver resolves to.
     */
    function pubkey() public view returns (string memory) {
        return _rsaPublicKey;
    }
}
