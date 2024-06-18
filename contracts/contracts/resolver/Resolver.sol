// SPDX-License-Identifier: MIT
pragma solidity >=0.8.4;

import "./resolvers/AddrResolver.sol";
import "./resolvers/RsaPublicKeyResolver.sol";
import "./resolvers/ServerResolver.sol";
import "./IBaseResolver.sol";

contract Resolver is IBaseResolver, AddrResolver, RsaPublicKeyResolver, ServerResolver {
    constructor(
        string memory rsaPublicKey,
        string memory serverEndpoint
    )
        AddrResolver(msg.sender)
        RsaPublicKeyResolver(rsaPublicKey)
        ServerResolver(serverEndpoint)
    {}
}
