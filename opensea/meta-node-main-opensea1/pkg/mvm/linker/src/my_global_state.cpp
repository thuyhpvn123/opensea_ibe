#include "my_global_state.h"
#include "mvm_linker.hpp"
#include "mvm/exception.h"
#include "mvm/gas.h"

struct GlobalStateGet_return {
  int status;
  unsigned char* balance_p;
  unsigned char* code_p;
  int code_length;
  unsigned char* storage_p;
  int storage_length;
};

namespace mvm
{
  using ET = Exception::Type;
  void MyGlobalState::remove(const Address& addr)
  {
    accounts.erase(addr);
  }

  AccountState MyGlobalState::get(const Address& addr, GasTracker *gas_tracker)
  {
    uint8_t *b_address = new uint8_t[32];
    mvm::to_big_endian(addr, b_address);

    const auto acc = accounts.find(addr);
    if (acc != accounts.cend()){
      if (gas_tracker != NULL) {
        gas_tracker->add_gas_used(getTouchedAddressGasCost());
      } 
      return acc->second;
    } else {
      GlobalStateGet_return accountQueryData = GlobalStateGet(b_address+12);
      if(accountQueryData.status == 2) {
        throw Exception(
                    ET::addressNotInRelated, 
                    "Address not in related addresses: " + mvm::address_to_hex_string(addr));
      }
      if(accountQueryData.status == 1) {
        uint8_t *copyBalance = new uint8_t[32];
        uint8_t *copyCode = new uint8_t[accountQueryData.code_length];
        uint8_t **copyStorage = new uint8_t*[accountQueryData.storage_length];

        memcpy(copyBalance, accountQueryData.balance_p, 32u);
        memcpy(copyCode, accountQueryData.code_p, accountQueryData.code_length);
        for (int i =0; i< accountQueryData.storage_length; i ++) {
          uint8_t *storage = new uint8_t[64];
          memcpy(storage, accountQueryData.storage_p + (i *64), 64);
          copyStorage[i] = storage; 
        }

        uint256_t balance = from_big_endian(copyBalance, 32u);
        std::vector<uint8_t> code(copyCode, copyCode+ accountQueryData.code_length);
      
        // accounts.erase(addr);
        insert({MyAccount(addr, balance, code), MyStorage(copyStorage, accountQueryData.storage_length)});
        const auto acc = accounts.find(addr);
        if (gas_tracker != NULL) {
          gas_tracker->add_gas_used(getUnTouchedAddressGasCost());
        } 
        return acc->second;
      }
      return create(addr, 0, {});
    }
  }

  AccountState MyGlobalState::create(
    const Address& addr, const uint256_t& balance, const Code& code)
  {
    insert({MyAccount(addr, balance, code), {}});
    
    return get(addr);
  }

  bool MyGlobalState::exists(const Address& addr)
  {
    return accounts.find(addr) != accounts.end();
  }

  size_t MyGlobalState::num_accounts()
  {
    return accounts.size();
  }

  const BlockContext& MyGlobalState::get_block_context()
  {
    return blockContext;
  }

  uint256_t MyGlobalState::get_block_hash(uint8_t offset)
  {
    return 0u;
  }

  uint256_t MyGlobalState::get_chain_id()
  {
    // TODO: may load from config
    return 0u;
  }

  void MyGlobalState::insert(const StateEntry& p)
  {
    const auto ib = accounts.insert(std::make_pair(p.first.get_address(), p));

    assert(ib.second);
  }

  bool operator==(const MyGlobalState& l, const MyGlobalState& r)
  {
    // TODO
    return true;
    // return (l.accounts == r.accounts) && (l.currentBlock == r.currentBlock);
  }

  // add changes functions
  void MyGlobalState::add_addresses_newly_deploy(const Address& addr, const Code& code) {
      addresses_newly_deploy[addr] = code;
  };

  void MyGlobalState::add_addresses_storage_change(const Address& addr, const uint256_t& key, const uint256_t& value) {
    cout << "Adding storage change: "<< key << ":" << value;

    addresses_storage_change[addr][key] = value;
  };

  void MyGlobalState::add_addresses_add_balance_change(const Address& addr, const uint256_t& amount) {
    addresses_add_balance_change[addr] += amount;
  };

  void MyGlobalState::add_addresses_sub_balance_change(const Address& addr, const uint256_t& amount) {
    addresses_sub_balance_change[addr] += amount;
  };

  uint8_t** MyGlobalState::get_newly_deploy(int& size, int* &code_sizes) {
    size = addresses_newly_deploy.size();
    uint8_t** rs = new uint8_t*[size];
    code_sizes = new int[size];
    int count = 0;
    for (const auto& p : addresses_newly_deploy)
    {
      int code_size = p.second.size(); 
      code_sizes[count] = code_size; 
      uint8_t* address_with_code = new uint8_t[32 + code_size];
      mvm::to_big_endian(p.first, address_with_code);
      std::memcpy(address_with_code + 32, p.second.data(), code_size);
      rs[count] = address_with_code;
      count ++;    
    }
    return rs;
  };

  uint8_t** MyGlobalState::get_storage_change(int& size, int* &storage_sizes) {
    size = addresses_storage_change.size();
    storage_sizes = new int[size];
    int total_storage_size = 0;
    int count = 0;
    for (const auto& p : addresses_storage_change)
    {
      int storage_size = 64 * p.second.size();
      storage_sizes[count] = storage_size;
      total_storage_size += 32 + storage_size; // 32 for address and storage size is total storage of that address
      count ++;
    }
    uint8_t** rs = new uint8_t*[total_storage_size];
    count = 0;
    for (const auto& p : addresses_storage_change)
    {
      uint8_t* storage = new uint8_t[storage_sizes[count]];
      int storage_count = 0;
      for (const auto& s : p.second) {
          int idx = storage_count * 64;
          cout << "Geting storage change: " << s.first << ":" << s.second;
          mvm::to_big_endian(s.first, storage + idx);
          mvm::to_big_endian(s.second, storage + idx + 32);
          storage_count++;
      }

      uint8_t* address_with_storage_change = new uint8_t[32 + storage_sizes[count]];
      mvm::to_big_endian(p.first, address_with_storage_change);
      std::memcpy(address_with_storage_change + 32, storage, storage_sizes[count]);
      rs[count] = address_with_storage_change;
      count ++;
    }
    return rs;
  };

  uint8_t** MyGlobalState::get_storage_root(int size) {
    uint8_t** rs = new uint8_t*[size]; // 32 bytes for address and 32 bytes for root hash
    int count = 0;
    for (const auto& p : addresses_storage_change)
    {
      uint8_t* address_with_storage_root = new uint8_t[64];
      mvm::to_big_endian(p.first, address_with_storage_root);
      MyStorage storage =  accounts[p.first].second;
      std::memcpy(address_with_storage_root + 32, storage.get_root_hash(), 32);
      rs[count] = address_with_storage_root;
      count ++;
    }

    return rs;
  }

  uint8_t** MyGlobalState::get_add_balance_change(int& size) {
    size = addresses_add_balance_change.size();
    uint8_t** rs = new uint8_t*[size];
    int count = 0;
    for (const auto& p : addresses_add_balance_change)
    {
      uint8_t* address_with_add_balance_change = new uint8_t[64];
      mvm::to_big_endian(p.first, address_with_add_balance_change);
      mvm::to_big_endian(p.second, address_with_add_balance_change + 32);
      rs[count] = address_with_add_balance_change;
      count++;
    }
    return rs;
  };

  uint8_t** MyGlobalState::get_sub_balance_change(int& size) {
    size = addresses_sub_balance_change.size();
    uint8_t** rs = new uint8_t*[size];
    int count = 0;
    for (const auto& p : addresses_sub_balance_change)
    {
      uint8_t* address_with_sub_balance_change = new uint8_t[64];
      mvm::to_big_endian(p.first, address_with_sub_balance_change);
      mvm::to_big_endian(p.second, address_with_sub_balance_change + 32);
      rs[count] = address_with_sub_balance_change;
      count++;
    }
    return rs;
  };
} // namespace mvm
