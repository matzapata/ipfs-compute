

TODO: Timeout docker running 15 secs or so
TODO: Pass method, headers and so on as args. Also allow more control on communication between processes so it works like a full server

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

## Contracts

### Credits escrow

```solidity

interface Escrow {
    // transfer usdc to escrow
    deposit(usdc amount) 

    withdraw()

    consume(price, request_hash, permit) onlyProvider -> (permit contains price, request hash)

    allowance(provider, budget) -> user allows provider to consume from his credits
}

```


### ENS registry

```solidity
map provider_domain -> resolver

register(domain, resolver)

getResolver(domain) -> resolver
```

### ENS resolver

```solidity
address()
public_key()
domain()
```

