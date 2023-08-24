pragma solidity >=0.4.22 <0.9.0;

contract TokenSwap {
    address public tokenContractAddress;
    uint public exchangeRate; // Fixed exchange rate between token and native coin

    constructor(address _tokenContractAddress, uint _exchangeRate) {
        tokenContractAddress = _tokenContractAddress;
        exchangeRate = _exchangeRate;
    }

    function swapTokensWithNativeCoin(uint tokenAmount) external payable {
        // Calculate the amount of native coin to be received based on the exchange rate
        uint nativeCoinAmount = tokenAmount * exchangeRate;

        // Verify that the caller has sent enough native coin
        require(msg.value >= nativeCoinAmount, "Insufficient native coin amount");

        // Perform the token transfer to the caller
        TokenContract tokenContract = TokenContract(tokenContractAddress);
        tokenContract.transfer(msg.sender, tokenAmount);

        // Send the native coin to the contract owner
        payable(owner()).transfer(nativeCoinAmount);
    }

    function swapNativeCoinWithTokens() external payable {
        // Calculate the amount of tokens to be received based on the exchange rate
        uint tokenAmount = msg.value / exchangeRate;

        // Perform the token transfer to the caller
        TokenContract tokenContract = TokenContract(tokenContractAddress);
        tokenContract.transfer(msg.sender, tokenAmount);
    }
}