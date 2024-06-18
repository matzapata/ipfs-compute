// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity  >=0.8.4;

contract Registry {
    struct Record {
        address owner;
        address resolver;
    }

    // Logged when the resolver for a node changes.
    event NewResolver(bytes32 indexed domain, address resolver);

    // each domain provider.eth is hashed -> bytes32
    mapping(bytes32 => Record) records;

    // Permits modifications only by the owner of the specified node.
    modifier onlyOwner(bytes32 domain) {
        address domainOwner = records[domain].owner;
        require(domainOwner == msg.sender);
        _;
    }

    /**
     * @dev Associates a domain with an owner and a resolver
     * @param domain The domain to update.
     * @param resolverAddres The address of the resolver.
     */
    function register(bytes32 domain, address resolverAddres) public {
        // Check if already registered
        require(records[domain].owner == address(0), "Domain is already registered");

        // Add new registry
        records[domain] = Record({
            owner: msg.sender,
            resolver: resolverAddres
        });

        // Emit the NewResolver event
        emit NewResolver(domain, resolverAddres);
    }

    /**
     * @dev Disassociates a domain with an owner and a resolver
     * @param domain The domain to update.
     */
    function unregister(bytes32 domain) public {
        // Check if already registered
        require(records[domain].owner == msg.sender, "Domain is not registered");

        // Remove registry
        delete records[domain];
    }

    
    /**
     * Returns the resolver associated with a domain.
     * @param domain The domain to query.
     */
    function resolver(
        bytes32 domain
    ) public view returns (address) {
        return records[domain].resolver;
    }


    /**
     * Returns the owner associated with a domain.
     * @param domain The domain to query.
     */
    function owner(
        bytes32 domain
    ) public view returns (address) {
        return records[domain].owner;
    }

    /**
     * @dev Updates the resolver associated with a domain.
     * @param domain The domain to update.
     * @param resolverAddres The address of the resolver.
     */
    function setResolver(
        bytes32 domain,
        address resolverAddres
    ) public onlyOwner(domain) {
        records[domain].resolver = resolverAddres;
        emit NewResolver(domain, resolverAddres);
    }
}