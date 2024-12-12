cmake_minimum_required(VERSION 3.25)
project(SDL_Test_2 C)
set(CMAKE_C_STANDARD 17)

add_executable(${PROJECT_NAME} WIN32 main.c UTF.c UTF.h std.c std.h console.h console.c SDL_is_stupid.h)

add_library(sdl SHARED IMPORTED)
set_target_properties(sdl PROPERTIES IMPORTED_IMPLIB "${CMAKE_SOURCE_DIR}/sdl/SDL2.lib")
set_target_properties(sdl PROPERTIES IMPORTED_LOCATION "${CMAKE_SOURCE_DIR}/sdl/SDL2.dll")
target_include_directories(${PROJECT_NAME} PUBLIC "${CMAKE_SOURCE_DIR}/sdl/include")

add_library(sdl_mixer SHARED IMPORTED)
set_target_properties(sdl_mixer PROPERTIES IMPORTED_IMPLIB "${CMAKE_SOURCE_DIR}/sdl_mixer/SDL2_mixer.lib")
set_target_properties(sdl_mixer PROPERTIES IMPORTED_LOCATION "${CMAKE_SOURCE_DIR}/sdl_mixer/SDL2_mixer.dll")
target_include_directories(${PROJECT_NAME} PUBLIC "${CMAKE_SOURCE_DIR}/sdl_mixer/include")
target_include_directories(${PROJECT_NAME} PUBLIC "${CMAKE_SOURCE_DIR}/sdl_mixer/stupid")

target_link_libraries(${PROJECT_NAME} sdl sdl_mixer)

# Debugging.
get_property(IncludeFolders TARGET ${PROJECT_NAME} PROPERTY INCLUDE_DIRECTORIES)
message(STATUS "IncludeFolders='${IncludeFolders}'")
