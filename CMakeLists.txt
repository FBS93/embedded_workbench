cmake_minimum_required(VERSION 3.22)
set(CMAKE_C_STANDARD 11) 
project(embedded_c_framework C ASM)

# Add custom cmake function to add executable.
include(tools/cmake/functions/ect_add_executable.cmake)

# Add src subdirectory
add_subdirectory(src)