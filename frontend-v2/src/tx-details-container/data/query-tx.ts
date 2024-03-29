import { gql } from "@apollo/client";

export const queryTx = gql`
  query QueryTx($hash: Buffer!) {
    transaction(hash: $hash) {
      id
      timestamp
      hash
      blockNumber
      value
      gasUsed
      cumulativeGasUsed
      error
      fromAddressHash
      toAddressHash
      status
      gas
      gasPrice
      index
      input
      nonce
      r
      s
      v
      revertReason
      maxPriorityFeePerGas
      maxFeePerGas
    }
  }
`;
