## BUILD
- pull all sub modules
```git submodule update --init --recursive```
```
export CGO_CFLAGS="-O -D__BLST_PORTABLE__"  
protoc --go_out=./pkg/ ./proto/*.proto -I ./proto/
go build .
```

## TODO
State channel:

## Native smart contract address 
0000000000000000000000000000000000000001
0000000000000000000000000000000000000002

## extension smart contract address
call api:     0000000000000000000000000000000000000101
extract json: 0000000000000000000000000000000000000102


//////////////////////////////////////////////////////// Offchain
## State Channel
1. Open state channel
2. Join state channel
3. Commit state channel
4. State channel node to execute transaction
5. Add transaction hash to commit data, sub transaction fee for state channel
6. Lock smart contract to state channel
7. Commit smart contract data to storage


/// split stake db
<!-- Add stake state root to block -->
Add udpate stake states to full block
Sync sender stake state in validator
sync receiver stake state in validator
sync sender stake state in node
sync receiver stake state in validator
update process transaction to stake and unstake
update schedule to load from stake state
update connections to load from stake state

test stake, test unstake
test get connection
test schedule

update db, use 1 db only for commit account states