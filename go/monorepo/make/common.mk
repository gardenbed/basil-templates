## PREREQUISITES:
##


## VARIABLES
##

nocolor := \033[0m
red := \033[1;31m
green := \033[1;32m
yellow := \033[1;33m
blue := \033[1;34m
purple := \033[1;35m
cyan := \033[1;36m


## MACROS
##

define echo_red
	echo "$(red)$(1)$(nocolor)"
endef

define echo_green
	echo "$(green)$(1)$(nocolor)"
endef

define echo_yellow
	echo "$(yellow)$(1)$(nocolor)"
endef

define echo_blue
	echo "$(blue)$(1)$(nocolor)"
endef

define echo_purple
	echo "$(purple)$(1)$(nocolor)"
endef

define echo_cyan
	echo "$(cyan)$(1)$(nocolor)"
endef
