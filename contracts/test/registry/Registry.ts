import { loadFixture } from "@nomicfoundation/hardhat-toolbox/network-helpers";
import { expect } from "chai";
import hre from "hardhat";


describe('Registry', function () {

    async function deployFixture() {
        const [owner, user] = await hre.ethers.getSigners();

        // Deploy Registry contract
        const Registry = await hre.ethers.getContractFactory('Registry');
        const registry = await Registry.deploy();

        // dummy values for testing
        const resolver = user.address; 
        const domain = hre.ethers.encodeBytes32String('provider.eth');
        
        return { owner, user, registry, resolver, domain };
    }

    it('should register a domain', async function () {
        const { owner, domain, registry, resolver } = await loadFixture(deployFixture);

        await registry.connect(owner).register(domain, resolver);

        expect(await registry.owner(domain)).to.equal(owner.address);
        expect(await registry.resolver(domain)).to.equal(resolver);
    });

    it('should not register a domain if already registered', async function () {
        const { owner, domain, registry, resolver } = await loadFixture(deployFixture);

        await registry.connect(owner).register(domain, resolver);

        await expect(
            registry.connect(owner).register(domain, resolver)
        ).to.be.revertedWith('Domain is already registered');
    });

    it('should unregister a domain', async function () {
        const { owner, domain, registry, resolver } = await loadFixture(deployFixture);


        await registry.connect(owner).register(domain, resolver);
        await registry.connect(owner).unregister(domain);

        expect(await registry.owner(domain)).to.equal(hre.ethers.ZeroAddress);
        expect(await registry.resolver(domain)).to.equal(hre.ethers.ZeroAddress);
    });

    it('should not unregister a domain if not owned', async function () {
        const { owner, user, domain, registry, resolver } = await loadFixture(deployFixture);

        await registry.connect(owner).register(domain, resolver);

        await expect(
            registry.connect(user).unregister(domain)
        ).to.be.revertedWith('Domain is not registered');
    });

    it('should update resolver for a domain', async function () {
        const { owner, domain, registry, resolver } = await loadFixture(deployFixture);

        const initialResolver = resolver
        const newResolver = owner.address; // dummy value for testing

        await registry.connect(owner).register(domain, initialResolver);
        await registry.connect(owner).setResolver(domain, newResolver);

        expect(await registry.resolver(domain)).to.equal(newResolver);
    });

    it('should not update resolver for a domain if not owned', async function () {
        const { owner, user, domain, registry, resolver } = await loadFixture(deployFixture);

        const initialResolver = resolver;
        const newResolver = owner.address; // dummy value for testing

        await registry.connect(owner).register(domain, initialResolver);

        await expect(
            registry.connect(user).setResolver(domain, newResolver)
        ).to.be.revertedWith('Only the owner can modify the resolver');
    });
});
