// SPDX-License-Identifier: MIT
pragma solidity >=0.8.4;

interface IServerResolver {
    function server() external view returns (string memory);
}
