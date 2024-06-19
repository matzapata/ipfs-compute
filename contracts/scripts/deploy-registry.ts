import hre from "hardhat";

async function main() {
    const [deployer] = await hre.ethers.getSigners();
    console.log("Deploying Registry with the account:", deployer.address);

    const Registry = await hre.ethers.getContractFactory("Registry");
    const registry = await Registry.deploy();
    console.log("Registry deployed to address:", await registry.getAddress());
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });
