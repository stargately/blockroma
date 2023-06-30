import Router from "koa-router";
import Koa from "koa";
import bodyParser from "koa-bodyparser";
import cors from "@koa/cors";
import {BigNumber} from "ethers";

(async function main() {
  const app = new Koa();
  const router = new Router();

  router.post("/", async (ctx) => {
    const body = ctx.request.body as any;

    console.log(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>\n" + JSON.stringify(body, null, 2));

    const resp = await ({

      "eth_chainId": async () => {
        return BigNumber.from(32382).toHexString()
      },


      "net_version": async () => {
        return "32382";
      },

      "eth_getBalance": async (_address: string, _height: string) => {
        return BigNumber.from("30000000000000000000").toHexString()
      },

      "eth_blockNumber": async () => {
        return "0x1759041"
      },

      "eth_getBlockByNumber": async () => {
        return {
          "difficulty": "0x4ea3f27bc",
          "extraData": "0x476574682f4c5649562f76312e302e302f6c696e75782f676f312e342e32",
          "gasLimit": "0x1388",
          "gasUsed": "0x0",
          "hash": "0xdc0818cf78f21a8e70579cb46a43643f78291264dda342ae31049421c82d21ae",
          "logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
          "miner": "0xbb7b8287f3f0a933474a79eae42cbca977791171",
          "mixHash": "0x4fffe9ae21f1c9e15207b1f472d5bbdd68c9595d461666602f2be20daf5e7843",
          "nonce": "0x689056015818adbe",
          "number": "0x1b4",
          "parentHash": "0xe99e022112df268087ea7eafaf4790497fd21dbeeb6bd7a1721df161a6657a54",
          "receiptsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
          "sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
          "size": "0x220",
          "stateRoot": "0xddc8b0234c2e0cad087c8b389aa7ef01f7d79b2570bccb77ce48648aa61c904d",
          "timestamp": "0x55ba467c",
          "totalDifficulty": "0x78ed983323d",
          "transactions": [],
          "transactionsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
          "uncles": []
        }
      },

      "eth_gasPrice": async () => {
        return "0x0";
      },

      "eth_estimateGas": async () => {
        return "0x5208";
      },

      "eth_getCode": async () => {
        return "0x1"; // should be the # of txs
      },

      "eth_getTransactionCount": async () => {
        return "0x1"; // should be the # of txs
      },

      "eth_call": async () => {
        return "0x5208"
      },

      "eth_sendRawTransaction": async (bytes: string) => {
        return "0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331";
      },

      "eth_getTransactionReceipt": async (txHash: string) => {
        return {
          "blockHash":
            "0xa957d47df264a31badc3ae823e10ac1d444b098d9b73d204c40426e57f47e8c3",
          "blockNumber": "0xeff35f",
          "contractAddress": null, // string of the address if it was created
          "cumulativeGasUsed": "0xa12515",
          "effectiveGasPrice": "0x5a9c688d4",
          "from": "0x6221a9c005f6e47eb398fd867784cacfdcfff4e7",
          "gasUsed": "0xb4c8",
          "logs": [{
            // logs as returned by getFilterLogs, etc.
          }],
          "logsBloom": "0x00...0", // 256 byte bloom filter
          "status": "0x1",
          "to": "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
          "transactionHash": txHash,
          "transactionIndex": "0x66",
          "type": "0x2"
        }
      },

      "eth_getBlockByHash": async (blockHash: string, _va: boolean) => {
        return {
          "difficulty": "0x4ea3f27bc",
          "extraData": "0x476574682f4c5649562f76312e302e302f6c696e75782f676f312e342e32",
          "gasLimit": "0x1388",
          "gasUsed": "0x0",
          "hash": blockHash,
          "logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
          "miner": "0xbb7b8287f3f0a933474a79eae42cbca977791171",
          "mixHash": "0x4fffe9ae21f1c9e15207b1f472d5bbdd68c9595d461666602f2be20daf5e7843",
          "nonce": "0x689056015818adbe",
          "number": "0x1b4",
          "parentHash": "0xe99e022112df268087ea7eafaf4790497fd21dbeeb6bd7a1721df161a6657a54",
          "receiptsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
          "sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
          "size": "0x220",
          "stateRoot": "0xddc8b0234c2e0cad087c8b389aa7ef01f7d79b2570bccb77ce48648aa61c904d",
          "timestamp": "0x55ba467c",
          "totalDifficulty": "0x78ed983323d",
          "transactions": [],
          "transactionsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
          "uncles": []
        };
      },

    } as Record<string, Function>)[body.method](...body.params)


    console.log("<<<<<<<<<<<<<<<<<<<<<<<<\n" + JSON.stringify(resp));

    ctx.response.body = {
      "id": body.id,
      "jsonrpc": "2.0",
      "result": resp
    }
  });

  app.use(bodyParser());
  app.use(cors())
  app.use(router.routes());
  const port = 8545;
  app.listen(port);



  console.log(`listening on http://localhost:${port}`);
})();
