// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

/// @title IVersioned
/// @notice An interface for semantically versioned contracts.
interface IVersioned {
    /// @notice Returns the semantic version of the contract
    function version() external pure returns (string memory _version);
}