# Parse token transfer

There are two places we could track token contracts

1. Token contract deployment
2. Token transfer

Our current solution is for token transfer.

## Models

- [x] TokenTransfer
- [ ] Token (contract)
  - [ ] ERC20
  - [ ] ERC721
  - [ ] ERC1155
- [ ] Nft Metadata

## Use Cases / Pages

- `/tokens` list all tokens
  - Token
  - Address
  - Total Supply
  - Holders Count
- `/token/:address` token details
  - Contract Address
  - Total Supply
  - Holders
  - Transfers
  - Token Type
- `/assets/:address/:tokenId` ID details (NFT)
  - Name + Id
  - Contract Address
  - NFT storage method
  - Number of transfers
  - Current holder
  - `contract.tokenURI(0)`
    returns [JSON](https://ikzttp.mypinata.cloud/ipfs/QmQFkLSQysj94s5GvTHPyzTxrawwtjgiiYS2TBLgrvw8CW/0) with fields:
    - name: string
    - image: string
    - attributes: Array<{trait_type: string, value}>
- `/address/:addressHash` contract address details
  - Token
  - Creator
  - Balance
  - Tokens
  - Transactions
  - Transfers
  - Gas Used
  - Last Balance Update
