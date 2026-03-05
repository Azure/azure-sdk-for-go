# Frontmatter Requirements

## Required Fields

```yaml
---
name: my-skill-name
description: A clear description of what this skill does and when to use it.
---
```

**`name` rules:**
- 1-64 characters
- Lowercase letters, numbers, and hyphens only (`a-z`, `0-9`, `-`)
- Must not start or end with `-`
- Must not contain consecutive hyphens (`--`)
- Must match the parent directory name

**`description` rules:**
- 1-1024 characters
- Should describe BOTH what the skill does AND when to use it
- Include keywords that help agents identify relevant tasks

## Optional Fields

```yaml
---
name: my-skill-name
description: Description here.
license: Apache-2.0
compatibility: Requires az CLI and docker
metadata:
  author: your-org
  version: "1.0"
---
```

## Activation Triggers

Use `WHEN:` with distinctive quoted trigger phrases (preferred for cross-model compatibility):

```yaml
# Good - WHEN: with quoted phrases (preferred)
description: "Perform Azure compliance assessments using azqr. WHEN: \"check compliance\", \"assess Azure resources\", \"run azqr\", \"review security posture\"."
```

```yaml
# Accepted - USE FOR: still works
description: "Perform Azure compliance assessments using azqr. USE FOR: \"check compliance\", \"assess Azure resources\", \"run azqr\"."
```

```yaml
# Bad - too vague
description: "Helps with Azure stuff."
```

> ⚠️ **Do NOT add "DO NOT USE FOR:" clauses.** They cause keyword contamination on Claude Sonnet and similar models. Use positive routing with distinctive `WHEN:` phrases instead.
