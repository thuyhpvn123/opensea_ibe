cmake_minimum_required(VERSION 3.10)

project(c_merkle_patricia_trie)

set(CMAKE_CXX_STANDARD 17)

if(MSVC)
  # Set Windows compiler options.
  add_compile_options(/W3 /std:c++latest)
else()
  # Set Linux compiler options
   # add_compile_options(-Wall -Werror)

  if("${CMAKE_CXX_COMPILER_ID}" STREQUAL "GNU")
    separate_arguments(COVERAGE_FLAGS UNIX_COMMAND "--coverage -fprofile-arcs -ftest-coverage")
    set(PLATFORM_SPECIFIC_TEST_LIBS "gcov")
  else()
    separate_arguments(COVERAGE_FLAGS UNIX_COMMAND "-fprofile-instr-generate -fcoverage-mapping")
    set(PLATFORM_SPECIFIC_TEST_LIBS "-fprofile-instr-generate")
  endif()
endif()

file(GLOB KECCAK_SOURCES
  3rdparty/keccak/*.c
)

add_subdirectory(3rdparty)

# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
# Options 
# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
option(RECORD_TRACE "Record a detailed trace of EVM execution during test runs" OFF)
if(RECORD_TRACE)
  add_definitions(-DRECORD_TRACE)
endif(RECORD_TRACE)

set(INCLUDE_DIRS
  ${CMAKE_CURRENT_SOURCE_DIR}/3rdparty
  ${CMAKE_CURRENT_SOURCE_DIR}/include
)

set(CORE_SRCS
  src/merkle_trie.cpp
)

# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
# Libraries
# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
add_library(c_merkle_patricia_trie STATIC
  ${CORE_SRCS}
  ${KECCAK_SOURCES}
)
target_include_directories(c_merkle_patricia_trie PRIVATE
  ${INCLUDE_DIRS}
)
target_link_libraries(c_merkle_patricia_trie
  intx::intx
)

install(TARGETS c_merkle_patricia_trie DESTINATION ./)