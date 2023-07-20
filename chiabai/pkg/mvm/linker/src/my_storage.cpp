// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

#include "my_storage.h"

#include "mvm/util.h"
#include "mvm/gas.h"

#include <ostream>

namespace mvm
{
  MyStorage::MyStorage()
  {
    s = new MerkleTrie(32,true,mvm::keccak_256);
  } 

  MyStorage::MyStorage(uint8_t** b_storage, int storage_count)
  {
    s = new MerkleTrie(32,true,mvm::keccak_256);
    for (int i = 0; i < storage_count; i++) {
      uint256_t key = mvm::from_big_endian((uint8_t *)b_storage[i]);
      uint256_t value = mvm::from_big_endian((uint8_t *)b_storage[i]+32);
       s->set(&key, &value);
    }
  }

  void MyStorage::store(const uint256_t& key, const uint256_t& value, GasTracker* gas_tracker)
  {

    if(gas_tracker != NULL) {
      uint256_t old_value = load(key);
      if(value == old_value) {
        gas_tracker->add_gas_used(getSstoreGasCost(old_value, value));
      }
    }
    
    s->set(&key, &value);
  }

  uint256_t MyStorage::load(const uint256_t& key, GasTracker* gas_tracker)
  { 
    auto node = s->get(key);
    if (gas_tracker != NULL) {
      // TODO: check touched storage
      gas_tracker->add_gas_used(getTouchedStorageGasCost());
    }

    if (node == NULL)
      return 0;
    return *node->pValue;
  }

  bool MyStorage::exists(const uint256_t& key)
  {
    return s->get(key) != NULL;
  }

  bool MyStorage::remove(const uint256_t& key)
  {
    if (s->remove(key)){
      return true;
    }
    return false;
  }


  uint8_t * MyStorage::get_root_hash()
  {
    return s->root->hash;
  }
  
  inline std::ostream& operator<<(std::ostream& os, const MyStorage& s)
  {
    // os << nlohmann::json(s).dump(2);
    //TODO
    return os;
  }

} // namespace mvm
