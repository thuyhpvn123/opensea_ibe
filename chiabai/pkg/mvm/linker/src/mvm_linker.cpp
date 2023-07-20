#include <iostream>
#include <tuple>
#include "mvm_linker.hpp"
#include <string>
#include "mvm/opcode.h"
#include "mvm/processor.h"
#include "my_global_state.h"
#include "merkle_trie.h"
#include "my_logger.h"
#include "my_extension.h"

#include <cassert>
#include <fmt/format_header_only.h>
#include <fstream>
#include <iostream>
#include <nlohmann/json.hpp>
#include <random>
#include <sstream>
#include <vector>
#include <fstream>

nlohmann::json vectorLogsHandlerToJson(mvm::VectorLogHandler logHandler)
{
  auto json_logs = nlohmann::json::array();
  for (const auto &log : logHandler.logs)
  {
    nlohmann::json json_log;
    mvm::to_json(json_log, log);
    json_logs.push_back(json_log);
  }

  return json_logs;
}

void append_argument(std::vector<uint8_t> &code, const uint256_t &arg)
{
  // To ABI encode a function call with a uint256_t (or Address) argument,
  // simply append the big-endian byte representation to the code (function
  // selector, or bin). ABI-encoding for more complicated types is more
  // complicated, so not shown in this sample.
  const auto pre_size = code.size();
  code.resize(pre_size + 32u);
  mvm::to_big_endian(arg, code.data() + pre_size);
}

// Run input as an EVM transaction, check the result and return the output
mvm::ExecResult run(
  mvm::MyGlobalState &gs,
  bool deploy,
  const mvm::Address &from,
  const mvm::Address &to,
  const uint256_t &amount,
  uint64_t gas_price,
  uint64_t gas_limit,
  mvm::VectorLogHandler &log_handler,
  const mvm::Code &input)
{
  mvm::Transaction tx(
    from,
    amount, 
    gas_price, 
    gas_limit
  );
  // Record a trace to aid debugging
  mvm::Trace tr;
  MyLogger logger = MyLogger();
  MyExtension extension = MyExtension();
  mvm::Processor p(gs, log_handler, extension, logger);
  // Run the transaction
  const auto exec_result = p.run(tx, deploy, from, gs.get(to), input, amount, &tr);
  return exec_result;
}

mvm::BlockContext CreateBlockContext(
  uint64_t prevrandao,
  uint64_t gas_limit, 
  uint64_t time,      
  uint64_t base_fee,
  uint256_t number,   
  uint256_t coinbase
)
{
  mvm::BlockContext block_context;
  block_context.prevrandao = prevrandao;
  block_context.gas_limit = gas_limit;
  block_context.time = time;
  block_context.base_fee = base_fee;
  block_context.number = number;
  block_context.coinbase = coinbase;
  return block_context;
}

ExecuteResult processResult (mvm::ExecResult result, mvm::MyGlobalState &gs, mvm::VectorLogHandler &log_handler) {
  // storage
  char *b_output =  (char *)malloc((int)result.output.size() * sizeof(char));
  int length_output = (int)result.output.size();

  int length_add_balance_change;
  uint8_t **b_add_balance_change = gs.get_add_balance_change(length_add_balance_change);

  int length_sub_balance_change;
  uint8_t **b_sub_balance_change = gs.get_sub_balance_change(length_sub_balance_change);

  int length_code_change;
  int *length_codes = NULL;
  uint8_t **b_code_change = gs.get_newly_deploy(length_code_change, length_codes);
  int length_storage_change;
  int *length_storages  = NULL;
  uint8_t **b_storage_change = gs.get_storage_change(length_storage_change, length_storages);

  uint8_t ** b_storage_roots = gs.get_storage_root(length_storage_change); 
  // logs
  auto json_logs = vectorLogsHandlerToJson(log_handler);
  std::string str_logs = json_logs.dump();

  ExecuteResult rs = ExecuteResult{
    b_exitReason : (char)result.er,
    b_exception : (char)result.ex,
    b_exmsg : (char *)malloc((int)result.exmsg.size() * sizeof(char)),
    length_exmsg : (int)result.exmsg.size(),

    b_output : b_output,
    length_output : length_output,
    b_add_balance_change : (char **)b_add_balance_change,
    length_add_balance_change : length_add_balance_change,

    b_sub_balance_change : (char **)b_sub_balance_change,
    length_sub_balance_change : length_sub_balance_change,

    b_code_change : (char **)b_code_change,
    length_code_change : length_code_change,
    length_codes : length_codes,
    
    b_storage_change : (char **)b_storage_change,
    length_storage_change : length_storage_change,
    length_storages : length_storages,
    b_storage_roots: (char **)b_storage_roots,

    b_logs : (char *)malloc((int)str_logs.size() * sizeof(char)),
    length_logs : (int)str_logs.size(),

    gas_used: result.gas_used
  };

  memcpy(rs.b_exmsg, (char *)result.exmsg.c_str(), result.exmsg.size());
  memcpy(rs.b_output, result.output.data(), result.output.size());
  memcpy(rs.b_logs, str_logs.c_str(), str_logs.size());
  return rs;
}

ExecuteResult deploy(
  // transaction data
  unsigned char* b_caller_address,
  unsigned char* b_caller_last_hash,
  unsigned char* b_contract_constructor,
  int contract_constructor_length,
  unsigned char* b_amount,
  unsigned long long gas_price,
  unsigned long long gas_limit,
  // block context data
  unsigned long long block_prevrandao,
  unsigned long long block_gas_limit,
  unsigned long long block_time,
  unsigned long long block_base_fee,
  unsigned char* b_block_number,
  unsigned char* b_block_coinbase
)
{

  // format argument to right data type
  uint256_t caller_address = mvm::from_big_endian((uint8_t *)b_caller_address, 20u);
  uint256_t caller_last_hash = mvm::from_big_endian((uint8_t *)b_caller_last_hash, 32u);
  std::vector<uint8_t> contract_constructor((uint8_t *)b_contract_constructor, (uint8_t *)b_contract_constructor + contract_constructor_length);
  uint256_t amount = mvm::from_big_endian((uint8_t *)b_amount, 32u);

  uint256_t block_number = mvm::from_big_endian((uint8_t *)b_block_number, 32u);
  uint256_t block_coinbase = mvm::from_big_endian((uint8_t *)b_block_coinbase, 20u);
  mvm::BlockContext blockContext = CreateBlockContext(
    block_prevrandao,
    block_gas_limit,
    block_time, 
    block_base_fee,
    block_number,   
    block_coinbase
  );
  mvm::MyGlobalState gs(blockContext);
  //  init env
  mvm::VectorLogHandler log_handler;
  const auto contract_address = mvm::generate_contract_address(caller_address, caller_last_hash);
  // Set this constructor as the contract's code body
  auto contract = gs.create(contract_address, 0u, contract_constructor);

  auto result = run(
    gs,
    true,
    caller_address,
    contract_address,
    amount,
    gas_price,
    gas_limit,
    log_handler,
    {});
  auto code = result.output;
  contract.acc.set_code(std::move(code));

  gs.add_addresses_newly_deploy(contract_address, code);
  // update output to contract address
  std::vector<uint8_t> deployed_address(32);

  mvm::to_big_endian(contract_address, deployed_address.data());
  std::vector<uint8_t> truncated_address(20);
  memcpy(truncated_address.data(), deployed_address.data()+12, 20);

  result.output = truncated_address;
  ExecuteResult rs = processResult(result, gs, log_handler);
  
  return rs;
}

ExecuteResult call(
  // transaction data
  unsigned char* b_caller_address,
  unsigned char* b_contract_address,
  unsigned char* b_input,
  int   length_input,
  unsigned char* b_amount,
  unsigned long long gas_price,
  unsigned long long gas_limit,
  // block context data
  unsigned long long block_prevrandao,
  unsigned long long block_gas_limit,
  unsigned long long block_time,
  unsigned long long block_base_fee,
  unsigned char* b_block_number,
  unsigned char* b_block_coinbase
)
{
  // format argument to right data type
  uint256_t caller_address = mvm::from_big_endian((uint8_t *)b_caller_address, 20u);
  uint256_t contract_address = mvm::from_big_endian((uint8_t *)b_contract_address, 20u);
  std::vector<uint8_t> input((uint8_t *)b_input, (uint8_t *)b_input + length_input);
  uint256_t amount = mvm::from_big_endian((uint8_t *)b_amount, 32u);

  uint256_t block_number = mvm::from_big_endian((uint8_t *)b_block_number, 32u);
  uint256_t block_coinbase = mvm::from_big_endian((uint8_t *)b_block_coinbase, 20u);

  mvm::BlockContext blockContext = CreateBlockContext(
    block_prevrandao,
    block_gas_limit,
    block_time, 
    block_base_fee,
    block_number,   
    block_coinbase
  );

  mvm::MyGlobalState gs(blockContext);
  //  init env
  mvm::VectorLogHandler log_handler;

  auto result = run(
    gs,
    false,
    caller_address,
    contract_address,
    amount,
    gas_price,
    gas_limit,
    log_handler,
    input
  );

  ExecuteResult rs = processResult(result, gs, log_handler);

  return rs;
}