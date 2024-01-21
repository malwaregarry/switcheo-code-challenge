The following diagram shows a high level overview of the components that interact with the broadcaster service. 

Additinally, an authorisation service is needed to check whether the client has admin privileges to retry a failed request. 
The authorisation service can be any authorisation service already in use by the organisation.

In the problem scenario, there are 3 different cases in which the broadcaster service can be invoked:
1. a client sending a broadcast request
2. an admin sending a retry request
3. a client requesting for a list of transactions

A simple way to handle the different request is to have separate endpoints for each case.

```golang
http.HandleFunc("/transaction/broadcast", handleBroadcastRequest)   // POST
http.HandleFunc("/transaction/retry", handleRetryRequest)           // POST
http.HandleFunc("/transaction", handleGetRequest)                   // GET
```
Each request will generate an asynchronous task, so that tasks can yield when it is blocked waiting for IO. 
No synchronisation is required as the blockchain is assumed to have its own synchronisation mechanism and most transactional database engines have synchronisation built in.

### Case 1: Client sending a broadcast request
After receiving a POST request from the client, the broadcaster will perform the broadcast in 4 steps. The steps are carried out sequentially, which means that if any step fails, the function will return.

```javascript
async handleBroadcastRequest(client, data) {
    // 1. Sign data from client
    signed_transaction = sign(data) 
    // 2. Broadcast signed transaction to blockchain (limit to ~30s)
    status = await broadcast_to_blockchain_30s(signed_transaction)
    // 3. Update state in database to facilitate (SUCCESS or FAILURE)
    await db_update(signed_transaction, status)
    // 4. Reply client with HTTP OK
    await reply_client(client, signed_transaction)              
}
```
Notably, the timeout for the broadcast RPC call is set to 30s so that we do not have transactions pending and taking up memory on the broadcaster for an arbitrarily long time. 
Furthermore, the broadcaster will wait for a successful broadcast before replying the client. This approach was taken as it greatly simplifies the case where the the broadcast service crashes during the handling of the request.

An alternative design, which can achieve faster response times and a higher success rate,
is to reply with `HTTP OK` before invoking the broadcast RPC and have a more advanced state management system to allow recovery from a crash.

### Case 2:  Admin  sending a retry request
This case is similar to the first one, except that the the authorisation level of the client needs to be checked, and the id of the transaction in the db is sent instead of data. 

```javascript
async handleRetryRequest(client, authorisation, txn_id) {
    // 1. Check authorisation level
    actual_auth = await get_auth(client)
    if (authorisation != actual_auth) {
        reply_client(client, error(NOT_AUTHORISED)) 
        return
    }
    // 2. Get transaction and check state
    signed_transaction, state = await db_get(txn_id)
    if (state != FAILURE) {
        reply_client(client, error(NOT_FAILURE)) 
        return
    }
    // 3. Broadcast signed transaction to blockchain (limit to ~30s)
    status = await broadcast_to_blockchain_30s(signed_transaction)      
    // 4. Reply client with HTTP OK
    await reply_client(client, signed_transaction)                     
}
```

### Case 3:  Client requesting for a list of transactions
In this case, the broadcaster service will simply retrieve the list of transactions and return it to the client. 
Before responding to the client, there is also an intermediate step to first process list contents.

```javascript
async handleGetTransactionsRequest(client) {
    // 1. Get list of transactions
    transaction_list = await db_get()
    // 2. Process transaction list
    processed_list = process(transaction_list)
    // 3. Reply client
    await reply_client(client, processed_list)                     
}
```

