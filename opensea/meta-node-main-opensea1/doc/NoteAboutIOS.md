## bug and how to fix
 - in hello2/finsdk/finsdk/miner/bls/libbls384_256.a(fp.o), building for iOS Simulator, but linking in object file built for iOS, file 'hello2/finsdk/finsdk/miner/bls/libbls384_256.a' for architecture arm64
- build proto
protoc -I ./proto/ --objc_out=finsdk/proto/proto_objc/ proto/*.proto 

+ add exclude architech arm64

- proto file arc error
+ add -fno-objc-arc to file in build phase


- remember when add lib then add path to search path in build setting

- common convert funcs
+ string to nsstring
NSString* result = [NSString stringWithUTF8String:param.c_str()];

+ bytes to nsdata
NSData *data = [NSData dataWithBytes:&theData length:1];

+ allocate mutable aray
[NSMutableArray arrayWithCapacity:length_code_change];

+ NSDATA to char[]
unsigned char b_caller_address[senderAddress.length];
[senderAddress getBytes:b_caller_address length:senderAddress.length];

+ bug in 
    if (sd == nil) {
        sd = [[SmartContractData alloc] init];
        sd.code = code;
    }

    should be
    if (sd == nil) {
        sd = [[SmartContractData alloc] init];
    }
    sd.code = code;