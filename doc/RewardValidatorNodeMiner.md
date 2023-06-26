# Reward validator

when validator become leader, it will take all profit from block


transaction will have default fee
opcode will have defaut gas price

bên mình chạy giá gas động hay tĩnh anh,
Như eth thì gas động, thằng nào trả nhiều hơn thì được ưu tiên xử lý trước cách này thì giúp tăng lợi nhuận cho hệ thống nhưng lúc cao điểm phí giao dịch sẽ cao
Giá gas tĩnh thì tạo sự công bằng, đảm bảo được cái tính giá giao dịch rẻ


validator => reward module => send 
transactionFee = baseFee + gasUseEstimate*GasPrice


reward module interface
{
    reward(amount)
    new(threshHold)
    addRewardReceiver(address,amount) {
        if amountOf(address)>threshHold {
            createTransactionToSendAmountToAddress(amountOf(address) - transactionfee, address)
        }
    }
    
}

// block will have totalTransactionFee in it, validator will split equal for all node that success verify that block

// when node success verify an transaction, execute an transaction it will save to a reward list, when node receive an reward transaction from parent then it will split it by percent for address in reward list
// when pass thresh hold, it will create transaction to send reward back to address paticipated

validator - node
validator dựa vào số lượng tiền stake để ra tỉ lệ làm leader
leader sẽ tạo block và nhận tiền phí của block

validator có danh sách node kết nối đến với nó, mỗi khi nhận được phí của block thì những node nào thành công tạo vote cho block sẽ được chia đều tiền phí của block (validator đã trích 1 phần trong đó)
số tiền này chưa được chuyển liền mà sẽ được lưu vào 1 map address => amount
mỗi khi amount vượt qua ngưỡng(được config) thì validator tiến hành tạo giao dịch chuyển tiền cho address đó


node - miner
mỗi khi có miner thực thi thành công, hoặc kiểm tra chữ ký thành công thì lưu lại trong 1 map
address => số lượng thực thi thành công
mỗi khi node nhận receipt giao dịch trả thưởng từ cha thì sẽ dựa vào map trên để chia lại % cho từng
address theo số lượng thực thi thành công (node tự giữ lại 1 ít cho bản thân)
cũng như validator là lưu vào 1 map address => amount
khi vượt ngưỡng thì toạ giao dịch chuyển tiền

unstake tốn 60 ngày tiền mới về
