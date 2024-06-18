// SPDX-License-Identifier: MIT
pragma solidity >=0.8.4;

interface IRsaPublicKeyResolver {
    function pubkey() external view returns (string memory);
}
