import { gql } from "@apollo/client";

export const queryTx = gql`
  query QueryTx($hash: Buffer!) {
    transaction(hash: $hash) {
      id
      timestamp
      hash
      blockNumber
      value
      valueWithDecimal
      gasUsed
      cumulativeGasUsed
      error
      fromAddressHash
      toAddressHash
      status
      gas
      gasPrice
      gasPriceWithDecimal
      index
      input
      nonce
      r
      s
      v
      revertReason
      maxPriorityFeePerGas
      maxPriorityFeePerGasWithDecimal
      maxFeePerGas
      maxFeePerGasWithDecimal
    }
  }
`;
