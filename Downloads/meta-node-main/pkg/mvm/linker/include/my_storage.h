// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

#pragma once

#include "mvm/storage.h"

#include <map>
#include <nlohmann/json.hpp>
#include "merkle_trie.h"

namespace mvm
{
  /**
   * merkle patricia trie implementation of Storage
   */
  class MyStorage : public Storage
  {
    MerkleTrie* s;
    
  public:
    MyStorage();
    MyStorage(uint8_t** b_storage, int storage_count);

    void store(const uint256_t& key, const uint256_t& value, GasTracker* gas_tracker = NULL) override;
    uint256_t load(const uint256_t& key, GasTracker* gas_tracker = NULL) override;
    bool remove(const uint256_t& key) override;
    bool exists(const uint256_t& key);
    uint8_t * get_root_hash();
  };
} // namespace mvm
