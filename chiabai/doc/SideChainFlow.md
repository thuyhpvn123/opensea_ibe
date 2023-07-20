# Vấn đề 1: 
nếu account có thể tạo được transaction trên nhiều nhánh validator thì khi 1 nhánh validator offline  
user thực hiện transaction trên nhánh đó rồi thực hiện trên các nhánh online thì hash giữa nhánh offline và online sẽ khác nhau
dẫn đến không thể đồng bộ sau này.

## Giải quyết: 
- Để chạy được thì trước tiên account sẽ cần phải fixed trên 1 validator, tức là transactions của account sẽ chỉ được xuất phát từ 1 validator.
Có thêm 1 loại thực thi để user update validator mà user muốn fixed trên đó. Điều kiện để update là cả 2 validator (validator hiện tại và validator  
user muốn update tới) đều phải online để tránh xung đột.
- smart contract cũng sẽ phải có fixed trên validator để có thể chạy offline

# Vấn đề 2: 
Khi chuyển tiền thì fromAddress phải là các account fixed trên validator để đảm bảo không bị double spending và đúng last hash
toAddress có thể là bất kỳ address nào do nhận tiền thôi thì không dính dáng gì double spending và nhận tiền cũng ko gây thay đổi last hash  
tuy nhiên accountStateRoot sẽ không đúng do accountStateRoot được tính từ cả pending balance
cmake
## Giải quyết: 
- Khi validator offline sẽ không còn tạo block, confirmBlock nữa mà tất cả giao dịch sẽ đi qua 1 quá trình xử lý nội bộ:
client => node => node => validator => node => votes from nodes => valid packs => update account states from transactions in packs => save transactions to offline processed pool

client gửi transaction lên node, node gửi lên các cấp trên đến validator
validator thay vì gửi đến next leader thì bây giờ tự xử, gửi về lại node để verify, thựck thi smart contract
khi có kết quả thì update lại account states và lưu transaction lại vào offline processed pool

Khi Validator online thì sẽ lấy hết transaction từ offline processed pool gửi cho next leader (hoặc đợi đến lượt làm leader),
leader sẽ ưu tiên các transaction từ offline processed pool trước, tiến hành thêm các transaction này vào block được tạo ra hiện tại.

# Vấn đề 3:
Thực thi smart contract khi offline sẽ gây thay đổi storage, lúc online thì dữ liệu storage đã bị thay đổi không thể tiến hành chạy lại để kiểm  
tra nữa

## Giải quyết: 
Có 2 cách để giải quyết ở chỗ này:
- hoặc là khi chạy offline mode thì storage phải lưu lại toàn bộ storage cuả từng bước thực hiện (cách này hơi tốn tài nguyên lưu trữ ở storage) để khi online các miner ở các nhánh  
khác có thể tiến hành chạy kiểm tra lại
- hoặc là sẽ làm kiểu private chain chỉ định 1 số miner nhất định được chạy 1 địa chỉ smart contract, trong từng transaction thực thi thì ở receipt sẽ đóng kèm chữ ký của miner thực thi lúc online chỉ  
cần kiểm tra transaction đó có đủ chữ ký miner chưa hay ko thôi (cần đánh giá thêm các trường hợp các smart contract khác nhau gọi nhau)
