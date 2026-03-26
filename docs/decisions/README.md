# Architecture Decision Records

This directory contains Architecture Decision Records (ADRs) for the dotfiles
repository.

## Index

| ADR | Title | Status | Date |
| ----- | ------- | -------- | ------ |
| [0001](0001-rewrite-tmux-icon-name-in-go.md) | Rewrite tmux-icon-name as a Go Bonzai command | Accepted | 2026-03-21 |
| [0002](0002-full-command-line-via-pane-pid.md) | Resolve full command line via pane_pid | Accepted | 2026-03-21 |
| [0003](0003-yaml-config-with-patterns-and-name-substitution.md) | YAML config with pattern matching and name substitution | Accepted | 2026-03-22 |
| [0004](0004-embedded-config-with-user-merge.md) | Embedded default config with user config merge | Accepted | 2026-03-22 |
| [0005](0005-flip-config-schema.md) | Change config schema: key = pattern, fields = name + icon | Accepted | 2026-03-22 |
| [0006](0006-full-setup-only.md) | Full Setup Only - No Tiered Configuration | Accepted | 2026-03-23 |
| [0007](0007-cli-command-tree-structure.md) | CLI Command Tree Structure | Proposed | - |
| [0008](0008-module-namespace-architecture.md) | Module and Namespace Architecture | Proposed | - |
| [0009](0009-init-command-behavior.md) | Init Command Behavior | Proposed | - |
| [0010](0010-x-utilities-not-managed-tools.md) | X Utilities Are Not Managed Tools | Proposed | - |
| [0011](0011-deps-recursive-tree-view.md) | Deps Command Shows Recursive Tree View | Proposed | - |

## Creating a New ADR

1. Copy the last ADR file as `NNNN-short-title.md`
2. Fill in the sections
3. Update this index
