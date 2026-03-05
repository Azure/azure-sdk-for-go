---
name: sensei
description: "**WORKFLOW SKILL** — Iteratively improve skill frontmatter compliance using the Ralph loop pattern. WHEN: \"run sensei\", \"sensei help\", \"improve skill\", \"fix frontmatter\", \"skill compliance\", \"frontmatter audit\", \"score skill\", \"check skill tokens\". INVOKES: token counting tools, test runners, git commands. FOR SINGLE OPERATIONS: use token CLI directly for counts/checks."
license: MIT
metadata:
  author: Microsoft
  version: "1.0.0"
---

# Sensei

> "A true master teaches not by telling, but by refining." - The Skill Sensei

Automates skill frontmatter improvement using the [Ralph loop pattern](https://github.com/soderlind/ralph) - iteratively improving skills until they reach Medium-High compliance with passing tests, then checking token usage and prompting for action.

## Help

When user says "sensei help" or asks how to use sensei, show this:

```
╔══════════════════════════════════════════════════════════════════╗
║  SENSEI - Skill Frontmatter Compliance Improver                  ║
╠══════════════════════════════════════════════════════════════════╣
║                                                                  ║
║  USAGE:                                                          ║
║    Run sensei on <skill-name>              # Single skill        ║
║    Run sensei on <skill-name> --skip-integration  # Fast mode    ║
║    Run sensei on <skill1>, <skill2>, ...   # Multiple skills     ║
║    Run sensei on all Low-adherence skills  # Batch by score      ║
║    Run sensei on all skills                # All skills       ║
║                                                                  ║
║  EXAMPLES:                                                       ║
║    Run sensei on appinsights-instrumentation                     ║
║    Run sensei on azure-security --skip-integration               ║
║    Run sensei on azure-security, azure-observability             ║
║    Run sensei on all Low-adherence skills                        ║
║                                                                  ║
║  WHAT IT DOES:                                                   ║
║    1. READ      - Load skill's SKILL.md, tests, and token count  ║
║    2. SCORE     - Check compliance (Low/Medium/Medium-High/High) ║
║    3. SCAFFOLD  - Create tests from template if missing          ║
║    4. IMPROVE   - Add WHEN: triggers (cross-model optimized)     ║
║    5. TEST      - Run tests, fix if needed                       ║
║    6. REFERENCES- Validate markdown links                        ║
║    7. TOKENS    - Check token budget, gather suggestions         ║
║    8. SUMMARY   - Show before/after with suggestions             ║
║    9. PROMPT    - Ask: Commit, Create Issue, or Skip?            ║
║   10. REPEAT    - Until Medium-High score + tests pass           ║
║                                                                  ║
║  TARGET SCORE: Medium-High                                       ║
║    ✓ Description > 150 chars, ≤ 60 words                         ║
║    ✓ Has "WHEN:" trigger phrases (preferred)                     ║
║    ✓ No "DO NOT USE FOR:" (risky in multi-skill envs)             ║
║    ✓ SKILL.md < 500 tokens (soft limit)                          ║
║                                                                  ║
║  MORE INFO:                                                      ║
║    See .github/skills/sensei/README.md for full documentation    ║
║                                                                  ║
╚══════════════════════════════════════════════════════════════════╝
```

## When to Use

- Improving a skill's frontmatter compliance score
- Adding trigger phrases and anti-triggers to skill descriptions
- Batch-improving multiple skills at once
- Auditing and fixing Low-adherence skills

## Invocation Modes

### Single Skill
```
Run sensei on azure-deploy
```

### Multiple Skills
```
Run sensei on azure-security, azure-observability
```

### By Adherence Level
```
Run sensei on all Low-adherence skills
```

### All Skills
```
Run sensei on all skills
```

## The Ralph Loop

For each skill, execute this loop until score >= Medium-High AND tests pass:

1. **READ** - Load `plugin/skills/{skill-name}/SKILL.md`, tests, and token count
2. **SCORE** - Run spec-based compliance check (see [SCORING.md](references/SCORING.md)):
   - Validate `name` per [agentskills.io spec](https://agentskills.io/specification) (no `--`, no start/end `-`, lowercase alphanumeric)
   - Check description length and word count (≤60 words)
   - Check triggers (WHEN: preferred, USE FOR: accepted)
   - Warn on "DO NOT USE FOR:" (risky in multi-skill environments)
   - Preserve optional spec fields (`license`, `metadata`, `allowed-tools`) if present
3. **CHECK** - If score >= Medium-High AND tests pass → go to TOKENS step
4. **SCAFFOLD** - If `tests/{skill-name}/` doesn't exist, create from `tests/_template/`
5. **IMPROVE FRONTMATTER** - Add WHEN: triggers (stay under 60 words and 1024 chars)
6. **IMPROVE TESTS** - Update `shouldTriggerPrompts` and `shouldNotTriggerPrompts` to match
7. **VERIFY** - Run `cd tests && npm test -- --testPathPattern={skill-name}`
8. **VALIDATE REFERENCES** - Run `cd scripts && npm run references {skill-name}` to check markdown links
9. **TOKENS** - Check token budget and line count (< 500 lines per spec), gather optimization suggestions
10. **SUMMARY** - Display before/after comparison with unimplemented suggestions
11. **PROMPT** - Ask user: Commit, Create Issue, or Skip?
12. **REPEAT** - Go to step 2 (max 5 iterations per skill)

## Scoring Criteria (Quick Reference)

Sensei validates skills against the [agentskills.io specification](https://agentskills.io/specification). See [SCORING.md](references/SCORING.md) for full details.

| Score | Requirements |
|-------|--------------|
| **Invalid** | Name fails spec validation (consecutive hyphens, start/end hyphen, uppercase, etc.) |
| **Low** | Basic description, no explicit triggers |
| **Medium** | Has trigger keywords/phrases, description > 150 chars, >60 words |
| **Medium-High** | Has "WHEN:" (preferred) or "USE FOR:" triggers, ≤60 words |
| **High** | Medium-High + compatibility field |

**Target: Medium-High** (distinctive triggers, concise description)

> ⚠️ "DO NOT USE FOR:" is **risky in multi-skill environments** (15+ overlapping skills) — causes keyword contamination on fast-pattern-matching models. Safe for small, isolated skill sets. Use positive routing with `WHEN:` for cross-model safety.

**Strongly recommended** (reported as suggestions if missing):
- `license` — identifies the license applied to the skill
- `metadata.version` — tracks the skill version for consumers

## Frontmatter Template

Per the [agentskills.io spec](https://agentskills.io/specification), required and optional fields:

```yaml
---
name: skill-name
description: "[ACTION VERB] [UNIQUE_DOMAIN]. [One clarifying sentence]. WHEN: \"trigger 1\", \"trigger 2\", \"trigger 3\"."
license: MIT
metadata:
  version: "1.0"
# Other optional spec fields — preserve if already present:
# metadata.author: example-org
# allowed-tools: Bash(git:*) Read
---
```

> **IMPORTANT:** Use inline double-quoted strings for descriptions. Do NOT use `>-` folded scalars (incompatible with skills.sh). Do NOT use `|` literal blocks (preserves newlines). Keep total description under 1024 characters and ≤60 words.

> ⚠️ **"DO NOT USE FOR:" carries context-dependent risk.** In multi-skill environments (10+ skills with overlapping domains), anti-trigger clauses introduce the very keywords that cause wrong-skill activation on Claude Sonnet and fast-pattern-matching models ([evidence](https://gist.github.com/kvenkatrajan/52e6e77f5560ca30640490b4cc65d109)). For small, isolated skill sets (1-5 skills), the risk is low. When in doubt, use positive routing with `WHEN:` and distinctive quoted phrases.

## Test Scaffolding

When tests don't exist, scaffold from `tests/_template/`:

```bash
cp -r tests/_template tests/{skill-name}
```

Then update:
1. `SKILL_NAME` constant in all test files
2. `shouldTriggerPrompts` - 5+ prompts matching new frontmatter triggers
3. `shouldNotTriggerPrompts` - 5+ prompts matching anti-triggers

**Commit Messages:**
```
sensei: improve {skill-name} frontmatter
```

## Constraints

- Only modify `plugin/skills/` - these are the Azure skills used by Copilot
- `.github/skills/` contains meta-skills like sensei for developer tooling
- Max 5 iterations per skill before moving on
- Description must stay under 1024 characters
- SKILL.md should stay under 500 tokens (soft limit)
- Tests must pass before prompting for action
- User chooses: Commit, Create Issue, or Skip after each skill

## Flags

| Flag | Description |
|------|-------------|
| `--skip-integration` | Skip integration tests for faster iteration. Only runs unit and trigger tests. |

> ⚠️ Skipping integration tests speeds up the loop but may miss runtime issues. Consider running full tests before final commit.

## Reference Documentation

- [SCORING.md](references/SCORING.md) - Detailed scoring criteria
- [LOOP.md](references/LOOP.md) - Ralph loop workflow details
- [EXAMPLES.md](references/EXAMPLES.md) - Before/after examples
- [TOKEN-INTEGRATION.md](references/TOKEN-INTEGRATION.md) - Token budget integration

## Related Skills

- [markdown-token-optimizer](/.github/skills/markdown-token-optimizer) - Token analysis and optimization
- [skill-authoring](/.github/skills/skill-authoring) - Skill writing guidelines
