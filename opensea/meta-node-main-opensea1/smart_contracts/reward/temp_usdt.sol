//SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;
import "token2/contracts/access/Ownable.sol";
import "token2/contracts/security/Pausable.sol";
import "token2/contracts/token/ERC20/ERC20.sol";
// import "./Pausable.sol";

error NotController(address _address);

contract VNDT is Ownable, ERC20, Pausable {
    // Mapping
    mapping(address => uint256) private _balances;
    mapping(address => mapping(address => uint256)) private _allowances;
    mapping(address => bool) public isBlackListed;
    mapping(address => bool) public controllers;

    // Global Variable
    uint256 private _totalSupply;
    uint8 private _decimals;
    string private _symbol;
    string private _name;
    address public VNDTController;

    // Event
    event DestroyedBlackFunds(address _blackListedUser, uint _balance);
    event AddedBlackList(address _user);
    event RemovedBlackList(address _user);
    event MintByController(
        address _controller,
        address _recipient,
        uint _amount
    );

    // address _trustedVNDTController
    constructor(
        address _rewardAddress
    ) ERC20("Meta Dollar Reward", "USDMR", _rewardAddress) {
        // VNDTController = _trustedVNDTController;
        // controllers[_trustedVNDTController] = true;
        _decimals = 18;
        _mint(owner(),100000);
    }

    // modifier onlyController() {
    //     require(controllers[msg.sender], "You're not the Controller");
    //     _;
    // }

    // function editController(
    //     address _controller,
    //     bool _status
    // ) external onlyOwner returns (bool) {
    //     controllers[_controller] = _status;
    //     return _status;
    // }

    function Pause() external onlyOwner {
        _pause();
    }

    function UnPause() external onlyOwner {
        _unpause();
    }

    function transfer(
        address to,
        uint256 amount
    ) public override whenNotPaused returns (bool) {
        _transfer(_msgSender(), to, amount);
        return true;
    }

    function transferFrom(
        address from,
        address to,
        uint256 amount
    ) public override whenNotPaused returns (bool) {
        bool success = super.transferFrom(from, to, amount);
        return success;
    }

    function burnByOwner(address account, uint256 amount) external onlyOwner {
        _burn(account, amount);
    }

    function getBlackListStatus(address _maker) external view returns (bool) {
        return isBlackListed[_maker];
    }

    function addBlackList(address _evilUser) public onlyOwner {
        isBlackListed[_evilUser] = true;
        emit AddedBlackList(_evilUser);
    }

    function removeBlackList(address _clearedUser) public onlyOwner {
        isBlackListed[_clearedUser] = false;
        emit RemovedBlackList(_clearedUser);
    }

    function destroyBlackFunds(address _blackListedUser) public onlyOwner {
        require(isBlackListed[_blackListedUser], "Not in blacklist");
        uint dirtyFunds = balanceOf(_blackListedUser);
        _balances[_blackListedUser] = 0;
        _totalSupply -= dirtyFunds;
        emit DestroyedBlackFunds(_blackListedUser, dirtyFunds);
    }

    function mintToAddress(
        address recipient,
        uint256 amount
    ) public onlyOwner returns (bool) {
        _mint(recipient, amount);
        return true;
    }

    // function mintByController(
    //     address recipient,
    //     uint amount
    // ) external onlyController returns (bool) {
    //     _mint(recipient, amount);
    //     emit MintByController(msg.sender, recipient, amount);
    //     return true;
    // }
}
