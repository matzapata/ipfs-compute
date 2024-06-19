import hre from "hardhat";

async function main() {
    // Parse the arguments
    const rsaPublicKey = process.env.RSA_PUBLIC_KEY;
    if (rsaPublicKey === undefined) {
        throw new Error("Please provide the RSA public key using environment variable RSA_PUBLIC_KEY");
    }
    const serverEndpoint = process.env.SERVER_ENDPOINT;
    if (serverEndpoint === undefined) {
        throw new Error("Please provide the server endpoint using environment variable SERVER_ENDPOINT");
    }

    const [deployer] = await hre.ethers.getSigners();
    console.log("Deploying Resolver with the account:", deployer.address);

    const Resolver = await hre.ethers.getContractFactory("Resolver");
    const resolver = await Resolver.deploy(rsaPublicKey, serverEndpoint);
    console.log("Resolver deployed to address:",await resolver.getAddress());
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });
