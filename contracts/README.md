# Deployed contracts

Escrow: 0x5Fe8861F6571174a9564365384AE9b01CcdCd8D6
Deployed with: 0x01aA5423d4671E27919b48d8023ab115559cbbe2

Registry: 0xdb42A86B1bfe04E75B2A5F2bF7a3BBB52D7FFD2F
Deployed with: 0x01aA5423d4671E27919b48d8023ab115559cbbe2

Resolver: 0xB4459f11f71E23e6E627694A702414017B7BfB96
Deployed with: 0x01aA5423d4671E27919b48d8023ab115559cbbe2

Resolver registered
Domain: provider.pelmeni
Owner: 0x01aA5423d4671E27919b48d8023ab115559cbbe2
Resolver: 0xB4459f11f71E23e6E627694A702414017B7BfB96


# Commands

## Scripts

Deploy escrow

```shell
export USDC_ADDRESS="0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359" 
npx hardhat run scripts/deploy-escrow.ts --network polygon
```

Deploy registry

```shell
npx hardhat run scripts/deploy-registry.ts --network polygon
```

Deploy resolver

```shell
export RSA_PUBLIC_KEY="-----BEGIN RSA PUBLIC KEY-----\nMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAwMA3Nf2qkqzUxFlXfcRn\nT84bcIOWvw4AjWio6o+PDV0hO9B0aKdQqb/2uks/LrQNP6fgjtAjtjf2CEdC5sfW\noBxJbyWLsxtrG05pXbKNED9Yr5ywhc19Q0LeE5dpJCXl1m2x5LUy+UkrpYquwn+7\n6fyMZXmf6SDJzrnxQADwUWPe8EEMtgupaF5B9eDwrcW0vXdq9+D/Ab4tDZUCoLOt\nlksHpd8vhDBAY+NAFcigr6qzAKkvuensvj0/ByahhhI4jp/uVeV7Hl3yMWn0HXgQ\nQpu/rR6lTe95yqyKOkvac2DEbwTLU72ePTgqLrKa4ycXgQNrdZJsJBeLHtNUrwWY\nXOMP0xLmhly5VN4NQTYGRZd+T58n0ArI3J/dch+ZxcRORKZVQiVJt7/O3pF5uMUe\nbYICRxzpmKqQvlRKKLb4pi93Lv3iHP+I0OQsr/g8/FNG9ArhHGxqzE/SR2fKBOoj\nZT2jtJb2JTM2zdbMWiORrXyHRQnPnFU9HAonj4MB8XdREZCoGxFuvXK/EbrUQL2F\nVSbbt40khKa0gC+lzNkC9ypVBSdnRhMn/gdT2JErZTm8kcehrRaHiMumiIEP8WGS\nEjVxTTT3mt3mo5DWCB0xQFVPHC3OctLSBjAEWCmARZtWNgea/5WTCUxuUCzWDME0\nnmvtvx9u/Y6f4ZTozOKj2C0CAwEAAQ==\n-----END RSA PUBLIC KEY-----"
export SERVER_ENDPOINT="localhost:4000"
npx hardhat run scripts/deploy-resolver.ts --network polygon 
```

Register domain

```shell
export RESOLVER_ADDRESS="0xB4459f11f71E23e6E627694A702414017B7BfB96"
export REGISTRY_ADDRESS="0xdb42A86B1bfe04E75B2A5F2bF7a3BBB52D7FFD2F"
export DOMAIN="provider.pelmeni"
npx hardhat run scripts/register-resolver.ts --network polygon
```

## Dev

```shell
npx hardhat help
npx hardhat test
REPORT_GAS=true npx hardhat test
npx hardhat node
npx hardhat ignition deploy ./ignition/modules/Lock.ts
```


