#!/bin/sh

have jira || return

eval "$(jira completion zsh)"
