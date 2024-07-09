
# Build

# Cli

## Setup

`make cli-install`

setup `.env` or globally define variables with:

```
ETH_RPC=""
IPFS_GATEWAY=""
IPFS_PINATA_APIKEY=""
IPFS_PINATA_SECRET=""
```

## Commands

### allowance

```
> khachapuri allowance provider.pelmeni --a 0x01aA5423d4671E27919b48d8023ab115559cbbe2

0x01aA5423d4671E27919b48d8023ab115559cbbe2 allowed provider.pelmeni to consume 1000000 USDC at 100 per request
```

### approve

```
> khachapuri approve provider.pelmeni 1000000 100 --pk 294545c494c2f772df6d566c9b3d4

approved provider.pelmeni to consume 1000000 USDC at 100 per request
tx: 0xecf9e7cbdbd31447a51ed6fa081b1fbd46bdc472c4bddbd751f2bdab8a790207
```

### deposit

```
khachapuri deposit 100000 --pk 294545c494c2f772df6d566c9b3d410

You are about to deposit 100000 USDC into the escrow account. Continue? (y/n): y

Deposit successful
approval hash: -
transaction hash: 0xaf2c947f47ec22043875538ebfefca6ebdb16829c1f7074c485b3afb12ad1cdd
```

### Withdrawal

```
> khachapuri withdraw 100000 --pk 294545c494c2f772df6d566c9b3d410329419df883490ca7f8d685b339d2837b

Withdrawal successful. 
Transaction hash: 0xd55f7261add57240109715bc8b202078eac86006dee95dafcf233faf0f562961
```

### balance

```
> khachapuri balance --a 0x01aA5423d4671E27919b48d8023ab115559cbbe2

Balance for address 0x01aA5423d4671E27919b48d8023ab115559cbbe2: 0
```

### resolve

```
> khachapuri resolve provider.pelmeni 

Provider address: 0x01aA5423d4671E27919b48d8023ab115559cbbe2
Provider public key: -----BEGIN RSA PUBLIC KEY-----\nMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAwMA3Nf2qkqzUxFlXfcRn\nT84bcIOWvw4AjWio6o+PDV0hO9B0aKdQqb/2uks/LrQNP6fgjtAjtjf2CEdC5sfW\noBxJbyWLsxtrG05pXbKNED9Yr5ywhc19Q0LeE5dpJCXl1m2x5LUy+UkrpYquwn+7\n6fyMZXmf6SDJzrnxQADwUWPe8EEMtgupaF5B9eDwrcW0vXdq9+D/Ab4tDZUCoLOt\nlksHpd8vhDBAY+NAFcigr6qzAKkvuensvj0/ByahhhI4jp/uVeV7Hl3yMWn0HXgQ\nQpu/rR6lTe95yqyKOkvac2DEbwTLU72ePTgqLrKa4ycXgQNrdZJsJBeLHtNUrwWY\nXOMP0xLmhly5VN4NQTYGRZd+T58n0ArI3J/dch+ZxcRORKZVQiVJt7/O3pF5uMUe\nbYICRxzpmKqQvlRKKLb4pi93Lv3iHP+I0OQsr/g8/FNG9ArhHGxqzE/SR2fKBOoj\nZT2jtJb2JTM2zdbMWiORrXyHRQnPnFU9HAonj4MB8XdREZCoGxFuvXK/EbrUQL2F\nVSbbt40khKa0gC+lzNkC9ypVBSdnRhMn/gdT2JErZTm8kcehrRaHiMumiIEP8WGS\nEjVxTTT3mt3mo5DWCB0xQFVPHC3OctLSBjAEWCmARZtWNgea/5WTCUxuUCzWDME0\nnmvtvx9u/Y6f4ZTozOKj2C0CAwEAAQ==\n-----END RSA PUBLIC KEY-----
Provider server endpoint: localhost:4000
```

# Provider

# Gateway

