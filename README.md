# imagen

A small Go CLI that writes a solid magenta PNG with its own size centered in black. **This repo was one-shotted by Claude Code. No human reviewed the code!**

```sh
mise install
mise run build
./bin/imagen -w 1024 -h 1024   # writes 1024x1024.png in the current directory
```

- `mise run test` runs the tests, `mise run lint` runs golangci-lint, `mise run check` runs both.
- The only third-party dependency is `golang.org/x/image` (the Go team's extended library), used for scalable text rendering with its bundled Go Bold font.
- An agent-agnostic [Agent Skill](https://agentskills.io) describing the tool lives in [`skills/imagen`](skills/imagen/SKILL.md).
