Buy item using any token


use token A to buy item that using token B

Get convert rate from A to mtd: swapA
Get convert rate from mtd to B: swapB


router smart contract:
Tracking smart contract with simple swap smart contract address


marketPlace:
ByWithToken(tokenAddress):
    from token address query swap smart contract address from router
    from accept token query swap smart contract address from router
    from price calculate mtd require
    from mtd require token
    swap token accept to mtd
    swap mtd to token 
    sub token to buy item


