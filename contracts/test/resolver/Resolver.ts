import { loadFixture } from "@nomicfoundation/hardhat-toolbox/network-helpers";
import { expect } from "chai";
import hre from "hardhat";


describe('Resolver', function () {
    const RSA_PUBLIC_KEY = '-----BEGIN RSA PRIVATE KEY-----\n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\n-----END RSA PRIVATE KEY-----';
    const SERVER_ENDPOINT = 'https://example.com';

    async function deployFixture() {
        const [owner] = await hre.ethers.getSigners();

        // Deploy Resolver contract
        const Resolver = await hre.ethers.getContractFactory('Resolver');
        const resolver = await Resolver.deploy(RSA_PUBLIC_KEY, SERVER_ENDPOINT);

        return { owner, resolver };
    }

    it('should resolve address correctly', async function () {
        const { owner, resolver } = await loadFixture(deployFixture);

        const resolvedAddr = await resolver.addr();
        expect(resolvedAddr).to.equal(owner.address);
    });

    it('should resolve RSA public key correctly', async function () {
        const { resolver } = await loadFixture(deployFixture);

        const resolvedPubkey = await resolver.pubkey();
        expect(resolvedPubkey).to.equal(RSA_PUBLIC_KEY);
    });

    it('should resolve server endpoint correctly', async function () {
        const { resolver } = await loadFixture(deployFixture);

        const resolvedServer = await resolver.server();
        expect(resolvedServer).to.equal(SERVER_ENDPOINT);
    });
});
