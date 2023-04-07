// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import { Claim, GameType, GameStatus, Timestamp } from "src/types/Types.sol";

import { IVersioned } from "src/interfaces/IVersioned.sol";
import { IBondManager } from "src/interfaces/IBondManager.sol";
import { IInitializable } from "src/interfaces/IInitializable.sol";

/// @title IDisputeGame
/// @author clabby <https://github.com/clabby>
/// @author refcell <https://github.com/refcell>
/// @notice The generic interface for a DisputeGame contract.
interface IDisputeGame is IInitializable, IVersioned {
    /// @notice Emitted when the game is resolved.
    /// TODO: Define the semantics of this event.
    event Resolved(GameStatus indexed status);

    /// @notice Returns the timestamp that the DisputeGame contract was created at.
    function createdAt() external view returns (Timestamp _createdAt);

    /// @notice Returns the current status of the game.
    function status() external view returns (GameStatus _status);

    /// @notice Getter for the game type.
    /// @dev `clones-with-immutable-args` argument #1
    /// @dev The reference impl should be entirely different depending on the type (fault, validity)
    ///      i.e. The game type should indicate the security model.
    /// @return _gameType The type of proof system being used.
    function gameType() external view returns (GameType _gameType);

    /// @notice Getter for the root claim.
    /// @return _rootClaim The root claim of the DisputeGame.
    /// @dev `clones-with-immutable-args` argument #2
    function rootClaim() external view returns (Claim _rootClaim);

    /// @notice Getter for the extra data.
    /// @dev `clones-with-immutable-args` argument #3
    /// @return _extraData Any extra data supplied to the dispute game contract by the creator.
    function extraData() external view returns (bytes memory _extraData);

    /// @notice Returns the address of the `BondManager` used 
    function bondManager() external view returns (IBondManager _bondManager);

    /// @notice If all necessary information has been gathered, this function should mark the game
    ///         status as either `CHALLENGER_WINS` or `DEFENDER_WINS` and return the status of
    ///         the resolved game. It is at this stage that the bonds should be awarded to the
    ///         necessary parties.
    /// @dev May only be called if the `status` is `IN_PROGRESS`.
    function resolve() external returns (GameStatus _status);
}