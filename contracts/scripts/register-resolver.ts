import hre from "hardhat";
import { keccak_256 as sha3 } from 'js-sha3';

async function main() {
    // Parse the arguments
    const resolverAddress = process.env.RESOLVER_ADDRESS;
    if (resolverAddress === undefined) {
        throw new Error("Please provide the Resolver address using environment variable RESOLVER_ADDRESS");
    }
    const domain = process.env.DOMAIN;
    if (domain === undefined) {
        throw new Error("Please provide the domain using using environment variable DOMAIN");
    }
    const registryAddress = process.env.REGISTRY_ADDRESS;
    if (registryAddress === undefined) {
        throw new Error("Please provide the address using environment variable REGISTRY_ADDRESS");
    }

    const [deployer] = await hre.ethers.getSigners();
    console.log("Registering Resolver with the account:", deployer.address);

    const Registry = await hre.ethers.getContractAt("Registry", registryAddress);
    const hashedDomain = "0x" + sha3(domain);
    const tx = await Registry.register(hashedDomain, resolverAddress);
    await tx.wait();

    // Check if the resolver is registered
    const owner = await Registry.owner(hashedDomain);
    const resolver = await Registry.resolver(hashedDomain);
    
    console.log("Resolver registered");
    console.log("Owner:", owner);
    console.log("Resolver:", resolver);
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });
