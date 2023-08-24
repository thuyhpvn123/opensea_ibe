# TODO

off chain


Tạo ra các offchain validator 
offchain validator sẽ là các validator xử lý offchain transaction 

Các validator này sẽ chịu trách nhiệm xử lý transaction cho các địa chỉ đã đồng ý tham gia mạng của validator đó
 

các state của các account nàu sẽ được xử lý ở trên validator này thôi


node sẽ có quyền commit state cuối cùng lên toàn mạng 


 
TODO:
    - refactor interface
    - refactor context
    - off chain: 
        + account sẽ tiến hành đăng ký vào 1 off chain node 
        + Khi account đã register vào off chain node thì amount trong transaction của user sẽ được chuyển cho off chain node 
        + node này sẽ có khả năng nhận và thực hiện giao dịch, trạng thái off chain sẽ được lưu lại trên node
        + node sẽ có lệnh commit để cập nhật lại trạng thái cho smart contract, refund balance cho account
        + smart contract khi deploy sẽ có thêm 1 lựa chọn cho off chain node
        + đối với các smart contract offchain, smart contract sẽ không được thực thi trực tiếp nữa mà phải thông qua offchain node để thực thi

     
    offchain node lúc này có thể xem như là private chain 