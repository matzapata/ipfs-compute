import { loadFixture } from "@nomicfoundation/hardhat-toolbox/network-helpers";
import { expect } from "chai";
import hre from "hardhat";

describe('Escrow', function () {
    const INITIAL_BALANCE = hre.ethers.parseUnits('1000', 6); // Initial USDC balance for testing (scaled to 6 decimal places)

    async function deployFixture() {
        // Contracts are deployed using the first signer/account by default
        const [owner, user1, user2, provider] = await hre.ethers.getSigners();

        // Deploy USDC token mock
        const UsdcToken = await hre.ethers.getContractFactory("UsdcMock");
        const usdcToken = await UsdcToken.deploy(owner.address);

        // Deploy escrow contract
        const Escrow = await hre.ethers.getContractFactory('Escrow');
        const escrow = await Escrow.deploy(usdcToken);

        // Mint initial USDC tokens to users for testing
        await usdcToken.mint(owner.address, INITIAL_BALANCE);
        await usdcToken.mint(user1.address, INITIAL_BALANCE);
        await usdcToken.mint(user2.address, INITIAL_BALANCE);

        return { owner, user1, user2, provider, usdcToken, escrow };
    }

    it('should deposit USDC into the escrow contract', async function () {
        const { user1, usdcToken, escrow } = await loadFixture(deployFixture);

        const depositAmount = hre.ethers.parseUnits('100', 6);

        await usdcToken.connect(user1).approve(escrow, depositAmount);
        await escrow.connect(user1).deposit(depositAmount);

        expect(await escrow.balanceOf(user1.address)).to.equal(depositAmount);
        expect(await usdcToken.balanceOf(escrow)).to.equal(depositAmount);
    });

    it('should withdraw USDC from the escrow contract', async function () {
        const { user1, usdcToken, escrow } = await loadFixture(deployFixture);

        // First deposit USDC into the escrow contract
        const depositAmount = hre.ethers.parseUnits('100', 6);
        await usdcToken.connect(user1).approve(escrow, depositAmount);
        await escrow.connect(user1).deposit(depositAmount);
        
        // Then withdraw USDC from the escrow contract
        const initialBalance = await usdcToken.balanceOf(user1.address);
        await escrow.connect(user1).withdraw(depositAmount);
        const finalBalance = await usdcToken.balanceOf(user1.address);
        expect(finalBalance - initialBalance).to.equal(depositAmount);
    });

    it('should allow provider to consume credits from user', async function () {
        const { user1, provider, usdcToken, escrow } = await loadFixture(deployFixture);

        const depositAmount = hre.ethers.parseUnits('100', 6);
        const price = hre.ethers.parseUnits('50', 6);

        // User deposits USDC and approves escrow
        await usdcToken.connect(user1).approve(escrow, depositAmount);
        await escrow.connect(user1).deposit(depositAmount);

        // Provider approves to consume credits from user
        await escrow.connect(user1).approve(provider.address, depositAmount, price);

        // Provider consumes credits from user
        await escrow.connect(provider).consume(user1.address, price);

        expect(await escrow.balanceOf(user1.address)).to.equal(depositAmount - price);
        expect(await escrow.balanceOf(provider.address)).to.equal(price);
    });

    it('should not allow provider to consume credits above approved price', async function () {
        const { user1, provider, usdcToken, escrow } = await loadFixture(deployFixture);

        const depositAmount = hre.ethers.parseUnits('100', 6);
        const price = hre.ethers.parseUnits('50', 6);
        const higherPrice = hre.ethers.parseUnits('60', 6);

        // User deposits USDC and approves escrow
        await usdcToken.connect(user1).approve(escrow, depositAmount);
        await escrow.connect(user1).deposit(depositAmount);

        // Provider approves to consume credits from user
        await escrow.connect(user1).approve(provider.address, depositAmount, price);

        // Provider tries to consume credits at a higher price
        await expect(
            escrow.connect(provider).consume(user1.address, higherPrice)
        ).to.be.revertedWith('Price is too high');
    });

    it('should not allow provider to consume credits above approved allowance', async function () {
        const { user1, provider, usdcToken, escrow } = await loadFixture(deployFixture);

        const depositAmount = hre.ethers.parseUnits('100', 6);
        const price = hre.ethers.parseUnits('50', 6);
        const higherAllowance = hre.ethers.parseUnits('40', 6);

        // User deposits USDC and approves escrow
        await usdcToken.connect(user1).approve(escrow, depositAmount);
        await escrow.connect(user1).deposit(depositAmount);

        // Provider approves to consume credits from user with a certain allowance
        await escrow.connect(user1).approve(provider.address, higherAllowance, price);

        // Provider tries to consume credits above the approved allowance
        await expect(
            escrow.connect(provider).consume(user1.address, price)
        ).to.be.revertedWith('Insufficient credits');
    });
});
