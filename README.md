# Khachapuri - Decentralized IPFS based compute

Imagine AWS lambdas for IPFS compute, host your code in IPFS and make it worldwide available. You have a server, offer your compute power and get paid, you want cheap decentralized serverless functions, host your code, protect your secrets and pay with usdc.

Why khachapuri? It's just mindbogglingly delicious and palta was already taken.

## Structure



## Design

Programs are binaries with curl interface. Imagine that, call a endpoint with curl but instead of processing the response elsewhere we process it right there, curl is the interface.


###  Allow provider spending

```mermaid
sequenceDiagram
    User->>+Escrow: Deposit usdc
    Escrow->>-User: ok
    User->>+Registry: resolve(hash(providerDomain))
    Registry->>-User: Resolver address
    User->>+Resolver: resolver.Addr()
    Resolver->>-User: Provider address
    User->>+Escrow: approve(providerAddress, pricePerRequest, budget)
    Escrow->>-User: ok
```

###  Make request to provider

The domain can also resolved by the gateway and that will be the preference

```mermaid
sequenceDiagram
    User->>+Registry: resolve(hash(providerDomain))
    Registry->>-User: Resolver address
    User->>+Resolver: resolver.Server()
    Resolver->>-User: Base url
    User->>+Provider: {method} {vaseUrl}/{cid} {data} {optional payer header}
    Provider->>-User: Response
```

###  Provider program execution

```mermaid
sequenceDiagram
    User->>+Provider: {method} {domain}/{cid} {data} {optional payer header (default to owner)}
    Provider->>+IPFS: GET {cid}
    IPFS->>-Provider: rsa(spec.json)
    Provider->>Provider: dec(rsa(spec.json))
    Provider->>+Escrow: Get allowance(payer)
    Escrow->>-Provider: 10 USDC available at 0.1 USDC per request
    alt Provider allowed to withdraw at price lower or equal than allowed
    Provider->>+Escrow: consume(price)
    Escrow->>Escrow: Transfer balance from requestor to provider
    Escrow->>-Provider: Ok
    Provider->>Provider: Process request
    Provider-->>User: Response
    else Provider now allowed due to budget or price
    Provider->>+Escrow: consume(price)
    Escrow->>-Provider: Fail
    Provider-->>-User: Fail
    end
```

###  Register provider

```mermaid
sequenceDiagram
    Provider->>+Resolver: deploy(address, domain, rsaPublicKey)
    Resolver->>-Provider: Resolver address
    Provider->>+Registry: register(hash(domain), resolverAddress)
    Registry->>-Provider: ok
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
    User->>+Registry: resolve(hash(providerDomain))
    Registry->>-User: Resolver address
    User->>+Resolver: resolver.RsaPubKey()
    Resolver->>-User: rsa public key
    User->>User: rsa(spec.json, provider public key)
    User->>+IPFS: encrypted spec
    IPFS->>-User: cid (callable)
```

### Gateway

We want to:
- Create an easy and secure entry point for the deployments.
- Allow for easy deployments
- Enable a reputation system
- Enable smart routing of requests

#### Reverse proxy

{provider}.khachapuri.xyz:  

Resolves {provider} and acts a universal entry safe point to providers

```mermaid
sequenceDiagram
    participant Requestor
    participant Gateway
    participant Registry as ProvidersRegistry
    participant ProviderServer

    Requestor->>+Gateway: Access https://{provider_domain}.khachapuri.xyz/{cid}
    Gateway->>+Registry: Resolve {provider_domain}
    Registry->>-Gateway: target
    Gateway->>+ProviderServer: Forward request to target
    ProviderServer->>-Gateway: Response
    Gateway->>-Requestor: Send response back to user
```

#### Local tunneling

Allow everybody to easily provide compute to the network. 

```mermaid
sequenceDiagram
    participant User as Requestor
    participant LTServer as Gateway
    participant Registry as ProvidersRegistry
    participant LTClient as Gateway Client
    participant LocalServer as Provider Server localhost:3000

    User->>+LTServer: Access https://{provider_domain}.khachapuri.xyz/{cid}
    LTServer->>+Registry: Resolve {provider_domain}
    Registry->>-LTServer: target
    LTServer->>+LTClient: Forward request via WebSocket
    LTClient->>+LocalServer: Forward request to local server
    LocalServer->>-LTClient: Send response back to LTClient
    LTClient->>-LTServer: Forward response via WebSocket
    LTServer->>-User: Send response back to user
```

## TODO

TODO: add to contract trial credits giveaway. Users cannot withdraw, neither providers, just provider assigns some free credits to some user for trial
TODO: reputation system
TODO: Timeout docker running 15 secs or so
TODO: Allow unencrypted deployments, anybody can run. Cool also if there can be more than one runner
TODO: option to encrypt deployment at rest
