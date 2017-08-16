# Passport Application Chaincode

## Application Overview

This application was written to demonstrate how money transfers can be modeled on the Blockchain.

This version of the chaincode uses Hyperledger Fabric v0.6.

## Available Chaincode APIs

The following chaincode APIs are available from both CLI and via JSON RPC:

### Deploy / Init APIs and Usage

*Usage (CLI)*

```
peer chaincode deploy -l golang -n mycc -c '{"Function": "Init", "Args":[""]}'
```

*Usage (JSON RPC)*

```
{
  "jsonrpc": "2.0",
  "method": "deploy",
  "params": {
    "type": 1,
    "chaincodeID": {
      "path": "https://github.com/mschimk1/passport-chaincode"
    },
    "ctorMsg": {
      "function": "init",
      "args": [
        ""
      ]
    },
    "secureContext": "user_type1_0"
  },
  "id": 1
}
```

### Invoke APIs and Usage

#### OpenAccount

  Opens an account. The account details are provided as a JSON string. A *customer_id* value must be provided.

*Usage (CLI)*

```
peer chaincode invoke -l golang -n mycc -c '{"Function": "OpenAccount", "Args":["{\"customer_id\":\"12345\", \"id\":\"1\", \"bank_name\":\"Test Bank\", \"account_holder\": \"Mike\", \"country\": \"AU\", \"currency\": \"AUD\", \"balance\":10000}"]}'
```

*Usage (JSON RPC)*

```
{
  "jsonrpc": "2.0",
  "method": "invoke",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "839c65506e38cc6b4f7193c628db8bb811b79fcb596ed681c83a44a05fbd5630002a7018cc1192cf762ed744278ae3e68ea091aa315cc7c119a5ca7dae132a48"
    },
    "ctorMsg": {
      "function": "OpenAccount",
      "args": [
        "{\"customer_id\":\"12345\", \"id\":\"1\", \"bank_name\":\"Test Bank\", \"account_holder\": \"Mike\", \"country\": \"AU\", \"currency\": \"AUD\", \"balance\":10000}"
      ]
    },
    "secureContext": "user_type1_0"
  },
  "id": 1
}
```

#### CloseAccount

*Usage (CLI)*

```
peer chaincode invoke -l golang -n mycc -c '{"Function": "CloseAccount", "Args":["12345", "1"]}'
```

*Usage (JSON RPC)*
```

```

#### TopupAccount

*Usage (CLI)*

```
peer chaincode invoke -l golang -n mycc -c '{"Function": "TopupAccount", "Args":["12345", "1", "9000"]}'
```

*Usage (JSON RPC)*
```
{
  "jsonrpc": "2.0",
  "method": "invoke",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "839c65506e38cc6b4f7193c628db8bb811b79fcb596ed681c83a44a05fbd5630002a7018cc1192cf762ed744278ae3e68ea091aa315cc7c119a5ca7dae132a48"
    },
    "ctorMsg": {
      "function": "TopupAccount",
      "args": [
        "12345", "1", "1100"
      ]
    },
    "secureContext": "user_type1_0"
  },
  "id": 1
}
```

#### TransferMoney

*Usage (CLI)*

```
peer chaincode invoke -l golang -n mycc -c '{"Function": "TransferMoney", "Args":["{\"from_customer\":\"1234\", \"from_account\":\"1\", \"to_customer\":\"5678\", \"to_account\":\"2\", \"currency\":\"AUD\", \"amount\":1000}"]}'
```

*Usage (JSON RPC)*
```
{
  "jsonrpc": "2.0",
  "method": "invoke",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "839c65506e38cc6b4f7193c628db8bb811b79fcb596ed681c83a44a05fbd5630002a7018cc1192cf762ed744278ae3e68ea091aa315cc7c119a5ca7dae132a48"
    },
    "ctorMsg": {
      "function": "TransferMoney",
      "args": [
        "{\"from_customer\":\"1234\", \"from_account\":\"1\", \"to_customer\":\"5678\", \"to_account\":\"2\", \"currency\":\"AUD\", \"amount\":1000}"
      ]
    },
    "secureContext": "user_type1_0"
  },
  "id": 1
}
```

### Query APIs and Usage

#### GetAccountList

*Usage (CLI)*

```
peer chaincode invoke -l golang -n mycc -c '{"Function": "GetAccountList", "Args":["12345"]}'
```

*Usage (JSON RPC)*
```
{
  "jsonrpc": "2.0",
  "method": "query",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "839c65506e38cc6b4f7193c628db8bb811b79fcb596ed681c83a44a05fbd5630002a7018cc1192cf762ed744278ae3e68ea091aa315cc7c119a5ca7dae132a48"
    },
    "ctorMsg": {
      "function": "GetAccountList",
      "args": [
        "12345"
      ]
    },
    "secureContext": "user_type1_0"
  },
  "id": 1
}
```

#### GetAccount

*Usage (CLI)*

```
peer chaincode invoke -l golang -n mycc -c '{"Function": "GetAccount", "Args":["1234", "1"]}'
```

*Usage (JSON RPC)*
```
{
  "jsonrpc": "2.0",
  "method": "query",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "839c65506e38cc6b4f7193c628db8bb811b79fcb596ed681c83a44a05fbd5630002a7018cc1192cf762ed744278ae3e68ea091aa315cc7c119a5ca7dae132a48"
    },
    "ctorMsg": {
      "function": "GetAccount",
      "args": [
        "1234", "1"
      ]
    },
    "secureContext": "user_type1_0"
  },
  "id": 1
}
```

#### GetTransactionList

*Usage (CLI)*

```
peer chaincode invoke -l golang -n mycc -c '{"Function": "GetTransactionList", "Args":["1234", "1"]}'
```

*Usage (JSON RPC)*
```
{
  "jsonrpc": "2.0",
  "method": "query",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "839c65506e38cc6b4f7193c628db8bb811b79fcb596ed681c83a44a05fbd5630002a7018cc1192cf762ed744278ae3e68ea091aa315cc7c119a5ca7dae132a48"
    },
    "ctorMsg": {
      "function": "GetTransactionList",
      "args": [
        "1234", "1"
      ]
    },
    "secureContext": "user_type1_0"
  },
  "id": 1
}
```

#### GetTransaction

*Usage (CLI)*

```
peer chaincode invoke -l golang -n mycc -c '{"Function": "GetTransaction", "Args":["1234", "1", "cc0f9b4d761e64e548827f2de4b49d8f"]}'
```

*Usage (JSON RPC)*
```
{
  "jsonrpc": "2.0",
  "method": "query",
  "params": {
    "type": 1,
    "chaincodeID": {
      "name": "839c65506e38cc6b4f7193c628db8bb811b79fcb596ed681c83a44a05fbd5630002a7018cc1192cf762ed744278ae3e68ea091aa315cc7c119a5ca7dae132a48"
    },
    "ctorMsg": {
      "function": "GetTransaction",
      "args": [
        "1234", "1", "cc0f9b4d761e64e548827f2de4b49d8f"
      ]
    },
    "secureContext": "user_type1_0"
  },
  "id": 1
}
```

## Notes

* This chaincode makes use of partial keys for account and transaction list queries

