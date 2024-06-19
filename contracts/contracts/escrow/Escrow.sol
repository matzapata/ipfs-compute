// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/utils/math/Math.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract Escrow is Ownable {
    using Math for uint256;

    IERC20 public usdcToken;

    // balance account has deposited
    mapping(address => uint256) public balances;

    // how much each user allows provider to consume
    mapping(address => mapping(address => uint256)) private _allowances;

    // how much each user accepts to pay to provider per request
    mapping(address => mapping(address => uint256)) private _prices;

    event Consumed(address provider, address user, uint256 amount);
    event Deposit(address user, uint256 amount);
    event Widthdraw(address user, uint256 amount);

    constructor(address _usdcTokenAddress) Ownable(msg.sender) {
        usdcToken = IERC20(_usdcTokenAddress);
    }

    /**
     * @dev Deposits USDC into the escrow contract
     * @param amount The amount of USDC to deposit
     */
    function deposit(uint256 amount) public {
        require(usdcToken.balanceOf(msg.sender) >= amount, "Insufficient USDC balance");
        require(usdcToken.allowance(msg.sender, address(this)) >= amount, "Allowance too low");

        bool success = usdcToken.transferFrom(msg.sender, address(this), amount);
        require(success, "USDC transfer failed");

        balances[msg.sender] += amount;

        emit Deposit(msg.sender, amount);
    }

    /**
     * @dev Withdraws USDC from the escrow contract
     * @param amount The amount of USDC to withdraw
     */
    function withdraw(uint256 amount) public {
        require(balances[msg.sender] >= amount, "Insufficient balance");
        balances[msg.sender] -= amount;

        bool success = usdcToken.transfer(msg.sender, amount);
        require(success, "USDC transfer failed");

        emit Widthdraw(msg.sender, amount);
    }

    /**
     * @dev Transfers balance from user to provider calling this function given the user allowed the provider at specified price
     * @param account The user to consume the credit from
     */
    function consume(address account, uint256 price) public {
        require(_prices[account][msg.sender] >= price, "Price is too high");
        require(_allowances[account][msg.sender] >= price, "Insufficient credits");
        
        _allowances[account][msg.sender] -= price;
        balances[account] -= price;
        balances[msg.sender] += price;

        emit Consumed(msg.sender, account, price);
    }

    /**
     * @dev Returns the balance of the specified account.
     */
    function balanceOf(address account) public view returns (uint256) {
        return balances[account];
    }

    /**
     * @dev Approve provider to consume credits from user at a certain price
     */
    function approve(address provider, uint256 amount, uint256 price) public {
        _allowances[msg.sender][provider] = amount;
        _prices[msg.sender][provider] = price;
    }


    /**
     * @dev Returns the allowance of the provider to consume credits from the user.
     */
    function allowance(address user, address provider) public view returns (uint256, uint256) {
        return (_allowances[user][provider], _prices[user][provider]);
    }
}
