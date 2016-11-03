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

```

#### CloseAccount

*Usage (CLI)*

```
peer chaincode invoke -l golang -n mycc -c '{"Function": "CloseAccount", "Args":["1"]}'
```

#### TopupAccount

*Usage (CLI)*

```
peer chaincode invoke -l golang -n mycc -c '{"Function": "TopupAccount", "Args":["1", "9000"]}'
```

#### TransferMoney

*Usage (CLI)*

```
peer chaincode invoke -l golang -n mycc -c '{"Function": "TransferMoney", "Args":["{\"from_account\":\"1\", \"to_account\":\"2\", \"currency\":\"AUD\", \"amount\":1000}"]}'
```

### Query APIs and Usage

#### GetAccountList

*Usage (CLI)*

```
peer chaincode invoke -l golang -n mycc -c '{"Function": "GetAccountList", "Args":["12345"]}'
```

#### GetAccount

*Usage (CLI)*

```
peer chaincode invoke -l golang -n mycc -c '{"Function": "GetAccount", "Args":["1"]}'
```

#### GetTransactionList

*Usage (CLI)*

```
peer chaincode invoke -l golang -n mycc -c '{"Function": "GetTransactionList", "Args":["1"]}'
```

#### GetTransaction

*Usage (CLI)*

```
peer chaincode invoke -l golang -n mycc -c '{"Function": "GetTransaction", "Args":["1", "47e1d9adcba83ca019c403db8ced444a9221849821be625e9edaffc6a791b119"]}'
```

## Notes

* Chaincode deploy / init method recreates the ledger data tables
* This version of the application assumes that bank account IDs are unique - it generates a random 8 digit ID if none is provided in an OpenAccount request


