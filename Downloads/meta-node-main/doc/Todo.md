# TODO

# Test

- multiple execute miner
- multiple everify miner
- multiple child node
- test tps
- test revert
- test restart miner
- test storage
- test explorer



TODO:

GAS price for OPCODES   V
TIME price for execute

fee = transaction fee + (gas use * gas fee) + tip fee + time price fee

add to transaction:
    fee
    maxGas
    tipFee
    maxTime


Update stake, unstake 

https://ethereum.org/en/developers/docs/evm/opcodes/

only 101 top staker become validator
update leader schedule 1hour
when unstake it will take 60 day(convert to block for easy calculate) to refund 
punnish and reward leader

node have to stake money for validator to join
miner have to stake money for node to join

Code interface for these funcs

commission flow 

user create transaction
commission sign

if have commission sign then smart contract will pay the fee

commision sign is smart contract creator sign

max_gas_use * gas_fee (1 + time_use)