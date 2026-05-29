---
name: skill-pipeline
description: "A pipeline that executes the entire process—from research to Skill/Subagent creation—with a single command. Simply specify a topic, and it automatically handles everything: Web research, best practice extraction, Skill/Subagent generation, and validation."
allowed-tools:
  - Task
  - Read
  - Write
  - Bash
model: sonnet
---

# Skill Pipeline

A pipeline that executes the entire process—from research to implementation—seamlessly from end to end.

## Prerequisites

**Required subagents** (`.claude/agents/`):
- `parallel-researcher.md` - Stage 1: Research
- `auto-validator.md` - Stage 4: Validate

**Optional subagents**:
- `skill-reviewer.md` - Optional quality check

**Error if missing**: Pipeline will report which subagent is unavailable and suggest creating it.

## Usage

```
/skill-pipeline [topic] [--type skill|subagent]
```

例:
- `/skill-pipeline "Go error handling best practices" --type skill`
- `/skill-pipeline "TypeScript Type Safety" --type subagent`

## Workflow

### Stage 1: Research (parallel-researcher)

```
Task: subagent_type=parallel-researcher
Prompt: "Investigated best practices regarding [topic]. Gathered the information necessary to create Skills and Subagents."
```

Output: `.research/[topic]/synthesis.json`

### Stage 2: Design

Generate Skill/Subagent Designs from Research Results:

```yaml
name: [derived from topic]
description: [1-line summary]
type: skill | subagent
triggers: [when to use]
core_functionality:
  - [capability 1]
  - [capability 2]
references_needed:
  - [reference doc 1]
tools_required:
  - [tool 1]
  - [tool 2]
```

### Stage 3: Generate

#### Skill の場合

```
.claude/skills/[name]/
├── SKILL.md           # Main Skill Definition
├── references/        # Reference Documents
│   └── best-practices.md
└── scripts/           # Scripts as needed
    └── validate.sh
```

#### Subagent の場合

```
.claude/agents/[name].md   # Sub-agent Definition
```

### Stage 4: Validate (auto-validator)

```
Task: subagent_type=auto-validator
Prompt: "Validate Generated Skills/Subagents: Syntax Check, Required Field Verification, and Best Practice Compliance"
```

### Stage 5: Report

Generate Final Report:

```json
{
  "pipeline_id": "skill-pipeline-[timestamp]",
  "topic": "Original topic",
  "type": "skill|subagent",
  "research_summary": "Key findings from research",
  "generated_files": [
    ".claude/skills/[name]/SKILL.md",
    ".claude/skills/[name]/references/best-practices.md"
  ],
  "validation_result": {
    "status": "pass|fail",
    "issues": []
  },
  "next_steps": [
    "I'll give it a test: /[skill-name]",
    "If there are any areas for improvement...: /skill-pipeline feedback"
  ]
}
```

## Pipeline Stages

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│  Research   │───▶│   Design    │───▶│  Generate   │───▶│  Validate   │───▶│   Report    │
│  (parallel) │    │             │    │             │    │  (auto)     │    │             │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
      │                  │                  │                  │                  │
      ▼                  ▼                  ▼                  ▼                  ▼
 .research/         design.yaml      .claude/skills/    validation.json    report.json
 synthesis.json                      or agents/
```

## Error Handling

| Stage | Failure | Recovery |
|-------|---------|----------|
| Research | WebSearch fails | Retry with alternative queries |
| Design | Ambiguous requirements | Ask user for clarification |
| Generate | Missing template | Use default template |
| Validate | Syntax errors | Auto-fix and retry |

## Options

| Option | Description | Default | Status |
|--------|-------------|---------|--------|
| `--type` | `skill` or `subagent` | Inferred | Implemented |
| `--skip-research` | Use existing research | `false` | Implemented |
| `--no-validate` | Skip validation | `false` | Implemented |
| `--output-dir` | Custom output directory | `.claude/` | Implemented |

## Execution Flow

1. **Parse input** - Extract topic and options
2. **Check existing research** - If `.research/[topic]/synthesis.json` exists and `--skip-research` not set, ask to reuse
3. **Spawn parallel-researcher** - If needed
4. **Wait for research** - Monitor agent completion
5. **Read synthesis** - Extract key findings
6. **Design structure** - Create skill/subagent spec
7. **Generate files** - Write SKILL.md or agent.md
8. **Spawn auto-validator** - Validate generated files
9. **Collect results** - Merge all outputs
10. **Report** - Display summary and next steps

## Example Session

```
User: /skill-pipeline "Claude Code hooks automation"

Pipeline Starting...

[Stage 1/5] Research
  ├─ Spawning parallel-researcher...
  ├─ 4 research agents deployed
  └─ Waiting for completion...

[Stage 2/5] Design
  ├─ Analyzing synthesis.json
  ├─ Type: skill (detected from topic)
  └─ Name: hooks-automation

[Stage 3/5] Generate
  ├─ Creating .claude/skills/hooks-automation/
  ├─ Writing SKILL.md
  └─ Writing references/hooks-patterns.md

[Stage 4/5] Validate
  ├─ Spawning auto-validator...
  ├─ Checking syntax: ✓
  ├─ Checking required fields: ✓
  └─ Checking best practices: ✓

[Stage 5/5] Report
  ✓ Pipeline completed successfully

  Generated:
  - .claude/skills/hooks-automation/SKILL.md
  - .claude/skills/hooks-automation/references/hooks-patterns.md

  Try it: /hooks-automation
```

## Integration with Other Skills

This skill orchestrates:
- `parallel-researcher` (subagent) - For comprehensive research
- `auto-validator` (subagent) - For validation
- `skill-reviewer` (subagent) - Optional quality check

## Best Practices Encoded

From Claude Code official guidelines:
- Skills are single-task focused
- Subagents handle multi-step workflows
- Always include `allowed-tools` or `tools`
- Use `model: haiku` for simple skills
- Include triggers in description
