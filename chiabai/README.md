## BUILD
- pull all sub modules
```git submodule update --init --recursive```
```
export CGO_CFLAGS="-O -D__BLST_PORTABLE__"  
protoc --go_out=./pkg/ ./proto/*.proto -I ./proto/
go build .
```

## TODO


## Native smart contract address 
0000000000000000000000000000000000000001
0000000000000000000000000000000000000002

## extension smart contract address
call api:     0000000000000000000000000000000000000101
extract json: 0000000000000000000000000000000000000102