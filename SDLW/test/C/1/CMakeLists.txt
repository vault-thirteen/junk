cmake_minimum_required(VERSION 3.25)
project(SDL_test_1 VERSION 1.0.0 DESCRIPTION "SDL Test #1" LANGUAGES C)
set_target_properties(PROPERTIES VERSION ${PROJECT_VERSION})
set(CMAKE_C_STANDARD 23)

add_executable(SDL_test_1 WIN32 main.c)

add_library(sdl SHARED IMPORTED)
set_target_properties(sdl PROPERTIES IMPORTED_IMPLIB "${CMAKE_SOURCE_DIR}/sdl/SDL2.lib")
set_target_properties(sdl PROPERTIES IMPORTED_LOCATION "${CMAKE_SOURCE_DIR}/sdl/SDL2.dll")
target_include_directories(${PROJECT_NAME} PUBLIC "${CMAKE_SOURCE_DIR}/sdl/include")

target_link_libraries(${PROJECT_NAME} sdl)

# Debugging.
get_property(IncludeFolders TARGET ${PROJECT_NAME} PROPERTY INCLUDE_DIRECTORIES)
message(STATUS "IncludeFolders='${IncludeFolders}'")
