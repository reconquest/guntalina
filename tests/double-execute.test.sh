#!/bin/bash

init

config <<DATA
actions:
    restart-software:
        commands:
            - echo executed

rules:
    - masks:
        - /etc/file1
      workflow:
        - reload-software
    - masks:
        - /etc/*
      workflow:
        - reload-software
DATA

sources <<DATA
/etc/file1
/etc/file2
DATA


tests_ensure run

tests_ensure assert_diff <<DATA
following actions will be executed:
reload-something
    echo command_1
    echo command_2

following commands will be executed:
echo command_1
echo command_2

executing: echo command_1
command_1

executing: echo command_2
command_2
DATA
