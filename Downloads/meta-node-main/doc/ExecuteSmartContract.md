#
1. client send transaction to node
2. node send transaction to verify miner to verify sign
3. node send transaction to parent node
4. repeat till reach validator
5. validator add to pool or forward to next leader
6. next leader create entry and send back to validator to verify 
7. validator send entry back to node to reverify
8. if have execute smart contract, node will send to execute miner or send to child node if full request



#
miner execute smart contract transaction
    - update code, storage to smart contract state
    - if revert then cancel
    - if success then commit
    - add new storage root to result

if finished all transaction then send result back to parent



update between 2 transaction
add to dirty

not yet save to db
add to pending


miner have account_states ? yes
    when process execute smart contract, it will copy account state and update on copy version
    if not revert then add that account state to account states
    if revert then not add
    both case create execute result? how
    save update info 
