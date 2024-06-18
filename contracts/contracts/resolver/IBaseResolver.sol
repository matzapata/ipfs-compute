// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./resolvers/IServerResolver.sol";
import "./resolvers/IRsaPublicKeyResolver.sol";
import "./resolvers/IAddrResolver.sol";

interface IBaseResolver is
    IServerResolver,
    IRsaPublicKeyResolver,
    IAddrResolver
{}
