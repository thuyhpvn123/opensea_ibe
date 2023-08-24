// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "token2/contracts/access/Ownable.sol";

interface IERC20 {
    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(
        address indexed owner,
        address indexed spender,
        uint256 value
    );

    function totalSupply() external view returns (uint256);

    function balanceOf(address account) external view returns (uint256);

    function transfer(address to, uint256 amount) external returns (bool);

    function allowance(
        address owner,
        address spender
    ) external view returns (uint256);

    function approve(address spender, uint256 amount) external returns (bool);

    function transferFrom(
        address from,
        address to,
        uint256 amount
    ) external returns (bool);
}

contract MinterSmC is Ownable {
    // Constant
    uint256 baseClaim = 1000; // MTD
    uint256[] private userHalving;
    uint256 private rateHalving;
    uint256 pow = 10 ** 18;
    uint256 public joinedAddrNumber;

    struct UserData {
        // Struct
        uint256 joinDate;
        uint256 balance;
        address referralToAddr;
    }

    struct DepositHistory {
        uint256 timestamp;
        uint256 rate1;
        uint256 usdtAmount;
        uint256 mtdAmount;
        bool isInitialzed;
    }

    mapping(address => UserData) public userDataStorage;
    mapping(address => bool) public isJoin;
    mapping(address => bool) public isReferral;
    mapping(address => uint256) public lastestMint;
    mapping(address => uint256) public claimVelSto; // claimVelocityStorage
    mapping(address => address[]) public securityCircle;
    mapping(address => address[]) public referralList;
    uint256 public inviIdCount;
    uint256 public referralIdCount;
    mapping(uint256 => Invitation) public inviStorage;
    mapping(uint256 => ReferralOfferInfo) public offerStorage;
    mapping(address => mapping(bytes32 => DepositHistory)) userToHistory;
    mapping(address => bytes32[]) userToArray;

    IERC20 public usdtContract;

    constructor(
        address _usdtAddress,
        uint256[] memory _initUserHalving,
        uint256 _rateHalving
    ) {
        usdtContract = IERC20(_usdtAddress);
        userHalving = _initUserHalving;
        rateHalving = _rateHalving * pow;
    }

    event Received(address from, uint256 value);

    receive() external payable {
        emit Received(msg.sender, msg.value);
    }

    modifier onlyJoin() {
        require(isJoin[msg.sender], "This address has not joined");
        _;
    }

    modifier onlyJoinReceiver(address addr) {
        require(isJoin[addr], "This address of receiver has not joined");
        _;
    }

    modifier onlyNotJoin() {
        require(!isJoin[msg.sender], "This address has already joined");
        _;
    }

    modifier onlyNotReferral(address _addr) {
        require(
            !isReferral[_addr],
            "This address of receiver has being referral"
        );
        _;
    }
    function setErc(address add) external onlyOwner {
        usdtContract = IERC20(add);
    }

    event JoinEvent(address indexed userAddress, uint256 indexed atMoment);

    function Join() public onlyNotJoin {
        isJoin[msg.sender] = true;
        userDataStorage[msg.sender] = UserData(block.timestamp, 0, address(0));
        claimVelSto[msg.sender] = baseClaim / 2;
        joinedAddrNumber++;
        UpdatePrice();
        emit JoinEvent(msg.sender, block.timestamp);
    }

    event Halving(uint256 rate);

    function halving() private {
        rateHalving = rateHalving / 2;
        emit Halving(rateHalving);
    }

    function removeElement(uint256 index) private {
        require(index < userHalving.length, "Invalid index");
        for (uint256 i = index; i < userHalving.length - 1; i++) {
            userHalving[i] = userHalving[i + 1];
        }
        userHalving.pop();
    }

    function UpdatePrice() private {
        if (userHalving.length == 0) {
            return;
        }
        if (joinedAddrNumber >= userHalving[0]) {
            halving();
            removeElement(0);
        }
    }

    event ClaimEvent(address indexed userAddress, uint256 claimedAmount);

    function Claim() public onlyJoin {
        require(
            block.timestamp - lastestMint[msg.sender] >= 1 hours,
            "Time for next claim is not enough"
        );
        uint referralClaim;
        for (uint i = 0; i < referralList[msg.sender].length; i++) 
        {       
            referralClaim = (referralList[msg.sender].length * claimVelSto[msg.sender]) / 10;
        }

        uint claimAmount = claimVelSto[msg.sender] + referralClaim;
        userDataStorage[msg.sender].balance += claimAmount;
        lastestMint[msg.sender] = block.timestamp;
        emit ClaimEvent(msg.sender, claimAmount);
    }

    event WithDraw(
        address withdrawer,
        uint256 amount,
        uint256 timestamp,
        bytes32 hash
    );

    function Withdraw() public onlyJoin {
        UserData memory userData = userDataStorage[msg.sender];
        uint256 usdtTransferAmount = (userData.balance * pow) / rateHalving;
        bool transferSuccess = usdtContract.transferFrom(
            msg.sender,
            address(this),
            usdtTransferAmount
        );
        require(transferSuccess, "Unsuccessful transfer USDT");
        (bool isSuccessful, ) = payable(msg.sender).call{
            value: (usdtTransferAmount * rateHalving) / pow
        }("");
        bytes32 historyhash = keccak256(
            abi.encodePacked(msg.sender, block.timestamp, usdtTransferAmount)
        );
        userToHistory[msg.sender][historyhash] = DepositHistory({
            timestamp: block.timestamp,
            rate1: rateHalving,
            usdtAmount: usdtTransferAmount,
            mtdAmount: (usdtTransferAmount * rateHalving) / pow,
            isInitialzed: true
        });

        require(isSuccessful, "Unsuccessful transfer MTD");
        userDataStorage[msg.sender].balance =
            userData.balance -
            ((usdtTransferAmount * rateHalving) / pow);
        emit WithDraw(
            msg.sender,
            usdtTransferAmount,
            block.timestamp,
            historyhash
        );
    }

    event CircleInvitation(
        uint256 indexed inviId,
        address indexed _from,
        address indexed _to,
        uint256 _atMoment
    );

    struct Invitation {
        address from;
        address to;
        bool acceptation;
    }

    function InviteToCircle(address _to) public onlyJoin onlyJoinReceiver(_to) {
        require(
            securityCircle[msg.sender].length < 5,
            "Exceed maximum number of security circle"
        );

        require(_to != msg.sender, "Cannot invite yourself");

        for (uint256 i = 0; i < securityCircle[msg.sender].length; i++) {
            require(
                _to != securityCircle[msg.sender][i],
                "Invitation receipent is aldready in the cirle"
            );
        }
        inviIdCount++;
        uint256 inviId = inviIdCount;
        inviStorage[inviId] = Invitation(msg.sender, _to, false);

        emit CircleInvitation(inviId, msg.sender, _to, block.timestamp);
    }

    event AcceptCircleInvitation(
        uint256 indexed inviId,
        address indexed _from,
        address indexed _to,
        uint256 _atMoment
    );

    // Chập nhận lệnh mời
    function AcceptInvitation(uint256 inviId) public onlyJoin {
        Invitation memory invitation = inviStorage[inviId];
        require(
            msg.sender == invitation.to,
            "Caller is not invitation receipent"
        );

        require(
            securityCircle[invitation.from].length < 5,
            "Exceed maximum number of security circle"
        );

        for (uint256 i = 0; i < securityCircle[invitation.from].length; i++) {
            require(
                invitation.to != securityCircle[invitation.from][i],
                "Invitation receipent is aldready in the cirle"
            );
        }

        inviStorage[inviId].acceptation = true;
        // Add new address into security circle of Inviter
        securityCircle[invitation.from].push(invitation.to);

        // Update base velocity of inviter
        uint256 increVeloClaim = baseClaim / 10;
        claimVelSto[invitation.from] += increVeloClaim;

        // Update veloClaim of address being referraled
        if (userDataStorage[invitation.from].referralToAddr != address(0)) {
            claimVelSto[userDataStorage[invitation.from].referralToAddr] +=
                increVeloClaim /
                10;
        }
        
        emit AcceptCircleInvitation(
            inviId,
            invitation.from,
            invitation.to,
            block.timestamp
        );
    }

    event Collected(address user, uint256 amountUSDT);

    function GiveBackDeposit(bytes32 hash) external payable onlyJoin {
        require(checkExist(hash), "Invalid hash history");
        DepositHistory storage history = userToHistory[msg.sender][hash];
        require(msg.value >= history.mtdAmount, "Not Enough Amount MTD");
        bool transferSuccess = usdtContract.transfer(
            msg.sender,
            history.usdtAmount
        );
        require(transferSuccess, "Unable to collect USDT");
        emit Collected(msg.sender, history.usdtAmount);
        delete userToHistory[msg.sender][hash];
    }

    function Deposit() public payable {} // Transfer money to contract

    function GenerateHistoryHash(
        uint256 timestamp,
        uint256 usdtAmount
    ) external view returns (bytes32) {
        return keccak256(abi.encodePacked(msg.sender, timestamp, usdtAmount));
    }

    function CheckHashInfo(
        bytes32 hash
    ) external view onlyJoin returns (DepositHistory memory) {
        return userToHistory[msg.sender][hash];
    }

    function moveOut(
        address payable _to,
        uint256 _amount
    ) external payable returns (bool) {
        require(_amount <= address(this).balance, "Invalid amount");
        (bool sent, ) = _to.call{value: _amount}("");
        require(sent, "Failed to send Ether");
        return sent;
    }

    function checkUserHalving()
        external
        view
        returns (uint256[] memory, uint256)
    {
        return (userHalving, userHalving.length);
    }

    function checkExist(bytes32 _hash) private view returns (bool) {
        DepositHistory memory history = userToHistory[msg.sender][_hash];
        if (history.isInitialzed) {
            return true;
        }
        return false;
    }

    event ReferralOffer(
        uint256 indexed offerId,
        address indexed _from,
        address indexed _to,
        uint256 _atMoment
    );

    struct ReferralOfferInfo {
        address from;
        address to;
        bool acceptation;
    }

    function OfferReferral(address _to) public onlyJoin onlyNotReferral(_to) {
        require(_to != msg.sender, "Cannot offer referral to yourself");

        require(
            referralList[msg.sender].length < 10,
            "Exceed maximum number of referral list"
        );

        referralIdCount++;
        uint256 referralId = referralIdCount;
        offerStorage[referralId] = ReferralOfferInfo(msg.sender, _to, false);

        emit ReferralOffer(referralId, msg.sender, _to, block.timestamp);
    }

    event ReferralOfferAcceptance(
        uint256 indexed referralId,
        address indexed _from,
        address indexed _to,
        uint256 _atMoment
    );

    function AcceptReferralOffer(
        uint256 referralId
    ) public onlyJoin onlyNotReferral(msg.sender) {
        ReferralOfferInfo memory offer = offerStorage[referralId];
        require(msg.sender == offer.to, "Caller is not invitation receipent");

        require(
            referralList[offer.from].length < 10,
            "Exceed maximum number of referral list"
        );

        inviStorage[referralId].acceptation = true;
        // Add new address into security circle of Inviter
        referralList[offer.from].push(offer.to);
        userDataStorage[offer.to].referralToAddr = offer.from;

        emit ReferralOfferAcceptance(
            referralId, 
            offer.from,
            offer.to,
            block.timestamp
        );
    }

    function GetReferralList(address _addr) public view returns(address[] memory){
        return referralList[_addr];
    }
}
