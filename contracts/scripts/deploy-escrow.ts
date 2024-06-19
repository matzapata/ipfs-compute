import hre from "hardhat";

async function main() {
    // Parse the arguments
    const usdcAddress = process.env.USDC_ADDRESS;
    if (usdcAddress === undefined) {
        throw new Error("Please provide the USDC address using environment variable USDC_ADDRESS");
    }

    const [deployer] = await hre.ethers.getSigners();
    console.log("Deploying Escrow with the account:", deployer.address);

    const Escrow = await hre.ethers.getContractFactory("Escrow");
    const escrow = await Escrow.deploy(usdcAddress);
    console.log("Escrow deployed to address:", await escrow.getAddress());
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });
