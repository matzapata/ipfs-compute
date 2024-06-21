# Khachapuri - Decentralized IPFS based compute

Imagine AWS lambdas for IPFS compute, host your code in IPFS and make it worldwide available. You have a server, offer your compute power and get paid, you want cheap decentralized serverless functions, host your code, protect your secrets and pay with usdc.

Why khachapuri? It's just mindbogglingly delicious and palta was already taken.


## Design

Programs are binaries with curl interface. Imagine that, call a endpoint with curl but instead of processing the response elsewhere we process it right there, curl is the interface.


###  Allow provider spending

```mermaid
sequenceDiagram
    User->>+Escrow: Deposit usdc
    Escrow->>-User: ok
    User->>+ENS: getResolver(providerDomain)
    ENS->>-User: resolverAddress
    User->>+Resolver: getAddress()
    Resolver->>-User: providerAddress
    User->>+Escrow: allow(providerAddress, pricePerCredit, credits)
    Escrow->>-User: ok
```

###  Make request to provider

```mermaid
sequenceDiagram
    User->>+ENS: getResolver(providerDomain)
    ENS->>-User: resolverAddress
    User->>+Resolver: getWebDomain()
    Resolver->>-User: domain
    User->>+Provider: POST {domain}/{cid} {data} {optional payer header}
    Provider->>-User: Response
```

###  Provider program execution

```mermaid
sequenceDiagram
    User->>+Provider: POST {domain}/{cid} {data} {optional payer header (default to owner)}
    Provider->>+IPFS: GET {cid}
    IPFS->>-Provider: spec.json
    Provider->>Provider: dec(spec.json)
    Provider->>+Escrow: Get allowance(payer)
    Escrow->>-Provider: 100 credits at 1 USDC per credit
    alt Provider allowed to withdraw 1 credit at price lower than allowed
    Provider->>+Escrow: consume(price, permit)
    Escrow->>Escrow: Reduce available credits
    Escrow->>-Provider: Ok
    Provider->>Provider: Process request
    Provider-->>User: Response
    else Provider now allowed due to credits or price
    Provider->>+Escrow: consume(price, permit)
    Escrow->>-Provider: FAIL
    Provider-->>-User: FAIL
    end
```

###  Register provider

```mermaid
sequenceDiagram
    Provider->>+Resolver: deploy(address, domain, publicKey)
    Resolver->>-Provider: resolverAddress
    Provider->>+ENS: register(hash(domain), resolverAddress)
    ENS->>-Provider: ok
```

###  Deployment

```mermaid
sequenceDiagram
    User->>User: Compile source code
    User->>+IPFS: Upload zipped program
    IPFS->>-User: cid
    User->>User: Attach cid to spec
    User->>User: Sign T&C and add to spec
    User->>User: Add env vars to spec
    User->>+ENS: getResolver(providerDomain)
    ENS->>-User: resolverAddress
    User->>+Resolver: getPublicKey()
    Resolver->>-User: PublicKey
    User->>User: Encrypt deployment spec
    User->>+IPFS: encrypted spec
    IPFS->>-User: deploymentCid
```


## TODOS

TODO: Timeout docker running 15 secs or so
TODO: Pass method, headers and so on as args. Also allow more control on communication between processes so it works like a full server
TODO: Create library for hashing domains (https://github.com/Arachnid/eth-ens-namehash/blob/master/index.js) maybe even assign all to owner and owner manages it (https://github.com/ensdomains/ens-contracts/blob/8e8cf71bc50fb1a5055dcf3d523d2ed54e725d28/contracts/registry/ENSRegistry.sol#L29)