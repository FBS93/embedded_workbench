#------------------------------------------------------------------------------
# Function: ect_add_executable
#
# Description:
#   Wrapper around add_executable() that automatically applies common project
#   settings to all executable targets.
#
# User-configurable global variables:
#   ECT_GLOBAL_LIBS     Libraries automatically linked to all executables.
#   ECT_STARTUP_FILES   Startup source files added to all executables.
#   ECT_LINKER_FILES    Linker scripts applied to all executables.
#
# This function:
#   - Calls add_executable() with the full argument list.
#   - Retrieves the target name (first argument).
#   - Links all libraries listed in ECT_GLOBAL_LIBS to the target.
#   - Adds all startup files listed in ECT_STARTUP_FILES as private sources.
#   - Applies all linker scripts listed in ECT_LINKER_FILES via -T options.
#   - Enables generation of a linker map file for the target.
#
# Parameters:
#   Same as add_executable(): first argument is the target name, followed by
#   the list of source files.
#------------------------------------------------------------------------------
function(ect_add_executable)
    # Call the real add_executable() with the full argument list.
    add_executable(${ARGV})

    # Target name is the first argument.
    set(target_name ${ARGV0})

    # Link all global libraries if any.
    foreach(lib IN LISTS ECT_GLOBAL_LIBS)
        target_link_libraries(${target_name} PRIVATE ${lib})
    endforeach()

    # Add all startup files if any.
    foreach(sf IN LISTS ECT_STARTUP_FILES)
        target_sources(${target_name} PRIVATE "${sf}")
    endforeach()
    
    # Add all linker scripts if any.
    foreach(ld IN LISTS ECT_LINKER_FILES)
        target_link_options(${target_name} PRIVATE "-T${ld}")
    endforeach()

    # Generate a .map file per target.
    target_link_options(${target_name} PRIVATE "-Wl,-Map=${target_name}.map")
endfunction()
