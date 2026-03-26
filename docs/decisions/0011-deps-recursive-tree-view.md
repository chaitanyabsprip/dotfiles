# ADR-0011: Deps Command Shows Recursive Tree View

Status: Proposed

## Context

The `dot deps` verb (ADR-0007) displays dependencies for managed tools. The presentation format and traversal depth need to be defined to give users a clear picture of what each tool depends on.

## Decision

`dot deps <tool>` displays a **recursive tree view** of the tool's dependencies.

### Behavior

| Command | Output |
|---------|--------|
| `dot deps tmux` | Recursive dependency tree for tmux |
| `dot deps` | Errors — requires a tool argument |

### Example Output

```
tmux
├── fzf (sessionizer)
├── tmux (binary)
└── ohmyposh
    └── unzip
```

### Rationale

- A tree view makes transitive dependencies visible at a glance.
- Recursive by default avoids users needing to manually chase dependency chains.
- Matches the mental model of `dot deps` as an inspection/diagnostic tool.

## Consequences

### Positive

- Users can see the full dependency chain before running `dot install`
- Helps diagnose missing dependencies on new machines
- Clear visual representation of tool relationships

### Negative

- Requires each tool to declare its dependencies in a structured way
- Deep dependency trees may produce verbose output
