# Passport Application Chaincode

## Application Overview

This application was written to demonstrate how money transfers can be modeled
on the Blockchain. Accounts and transactions are represented as ledger tables.

This version of the chaincode uses the Hyperledger Fabric v0.5-developer-preview / IBM Bluemix Blockchain Service v0.4.2.

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
peer chaincode invoke -l golang -n mycc -c '{"Function": "CloseAccount", "Args":["1"]}'
```

*Usage (JSON RPC)*
```

```

#### TopupAccount

*Usage (CLI)*

```
peer chaincode invoke -l golang -n mycc -c '{"Function": "TopupAccount", "Args":["1", "9000"]}'
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
        "1", "1100"
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
peer chaincode invoke -l golang -n mycc -c '{"Function": "TransferMoney", "Args":["{\"from_account\":\"1\", \"to_account\":\"2\", \"currency\":\"AUD\", \"amount\":1000}"]}'
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
        "{\"from_account\":\"1\", \"to_account\":\"2\", \"currency\":\"AUD\", \"amount\":1000}"
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
peer chaincode invoke -l golang -n mycc -c '{"Function": "GetAccount", "Args":["1"]}'
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
        "1"
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
peer chaincode invoke -l golang -n mycc -c '{"Function": "GetTransactionList", "Args":["1"]}'
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
        "1"
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
peer chaincode invoke -l golang -n mycc -c '{"Function": "GetTransaction", "Args":["1", "47e1d9adcba83ca019c403db8ced444a9221849821be625e9edaffc6a791b119"]}'
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
        "1", "6086c4be7c5bfe31c20fb445a48b379e239d9a8d0143a3ca691b32a50cf05615"
      ]
    },
    "secureContext": "user_type1_0"
  },
  "id": 1
}
```

## Notes

* Chaincode deploy / init method recreates the ledger data tables
* This version of the application assumes that bank account IDs are unique - it generates a random 8 digit ID if none is provided in an OpenAccount request


