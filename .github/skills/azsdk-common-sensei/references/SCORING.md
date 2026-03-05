# Scoring Criteria

Detailed scoring criteria for evaluating skill frontmatter compliance per the [agentskills.io specification](https://agentskills.io/specification).

## Overview

Sensei evaluates skills on two dimensions:
1. **Frontmatter Compliance** - Triggers, anti-triggers, description quality
2. **Token Budget** - Staying within recommended limits

## Token Budgets

From [skill-authoring](/.github/skills/skill-authoring):

| File | Soft Limit | Hard Limit | Notes |
|------|------------|------------|-------|
| SKILL.md | 500 | 5000 | Keep lean, use references for detail |
| references/*.md | 1000 | - | Each reference file |
| Description | - | 1024 chars | Frontmatter description field |

**Check with:** `cd scripts && npm run tokens -- check plugin/skills/{skill}/SKILL.md`

> **Units note:** Sensei measures in **tokens** (cl100k_base tokenizer), not words. Anthropic's [Complete Guide](https://resources.anthropic.com/hubfs/The-Complete-Guide-to-Building-Skill-for-Claude.pdf) recommends "under 5,000 words" for SKILL.md, while the [Agent Skills spec](https://agentskills.io/specification) recommends "< 5000 tokens" and "under 500 lines." Sensei uses the spec's token-based limits. As a rough conversion: 5000 tokens ≈ 3,750 words.

### Reference Loading Impact

References load **only when explicitly linked** in SKILL.md, not on activation:
- Use `[text](references/file.md)` to trigger loading
- Each file loads in full (entire content, not sections)
- No caching between requests
- Structure with recipes/services patterns for multi-option skills

See [skill-authoring REFERENCE-LOADING.md](/.github/skills/skill-authoring/references/REFERENCE-LOADING.md) for details.

## Adherence Levels

### Low Adherence

A skill is **Low** if:
- Description is < 150 characters (too brief to be useful)
- No explicit trigger phrases in description
- No anti-triggers
- Agent cannot reliably determine when to activate

**Examples of Low-adherence descriptions:**
```yaml
# Too brief (71 chars)
description: 'Instrument a webapp to send useful telemetry data to Azure App Insights'

# No triggers (just a catalog)
description: 'Azure Security Services including Key Vault, Managed Identity, RBAC, Entra ID, and Defender.'
```

### Medium Adherence

A skill is **Medium** if:
- Description > 150 characters
- Has implicit or explicit trigger keywords
- May have "TRIGGERS:" or "Use this skill when" language
- Still missing anti-triggers

**Examples of Medium-adherence descriptions:**
```yaml
description: >-
  Deploy applications to Azure using Azure Developer CLI (azd). USE THIS SKILL 
  when users want to deploy, publish, host, or run their application on Azure. 
  Trigger phrases include "deploy to Azure", "host on Azure", "publish to Azure".
```

### Medium-High Adherence (Target)

A skill is **Medium-High** if:
- Description > 150 characters and ≤ 60 words
- Has explicit trigger phrases via `WHEN:` (preferred) or `USE FOR:`
- Leads with unique action verb + domain in first sentence

> **Note:** `WHEN:` scores higher than `USE FOR:` because quoted trigger phrases are more distinctive for cross-model pattern matching.

> ⚠️ **"DO NOT USE FOR:" is risky in multi-skill environments.** Anti-trigger clauses introduce the very keywords that cause wrong-skill activation on Claude Sonnet and other models that use fast pattern matching rather than deep negation reasoning. Use positive routing instead. See Check 4 for context-dependent guidance.

**Example of Medium-High adherence (cross-model optimized):**
```yaml
description: "Instrument web apps to send telemetry to Azure Application Insights for monitoring and diagnostics. WHEN: \"add App Insights\", \"instrument my app\", \"set up monitoring\", \"add telemetry\"."
```

### High Adherence

A skill is **High** if:
- All Medium-High criteria met
- Has `compatibility` field documenting requirements
- Has examples section (optional but recommended)

**Strongly recommended** (reported as suggestions if missing):
- `license` field — identifies the license applied to the skill
- `metadata.version` — tracks the skill version for consumers

**Example of High adherence:**
```yaml
name: appinsights-instrumentation
description: "Instrument web apps to send telemetry to Azure Application Insights for monitoring and diagnostics. WHEN: \"add App Insights\", \"instrument my app\", \"set up monitoring\", \"add telemetry\", \"track requests\"."
license: MIT
compatibility: Supports ASP.NET Core (.NET 6+), Node.js. Requires App Insights resource.
metadata:
  author: example-org
  version: "1.0.0"
```

---

## Rule-Based Checks

### 1. Name Validation

Per the [agentskills.io spec](https://agentskills.io/specification), the `name` field:
- Must be 1-64 characters
- May only contain lowercase alphanumeric characters and hyphens (`a-z`, `0-9`, `-`)
- Must not start or end with a hyphen
- Must not contain consecutive hyphens (`--`)
- Must match the parent directory name

| Check | Pass | Fail |
|-------|------|------|
| Lowercase only | `azure-deploy` | `Azure-Deploy` |
| Alphanumeric + hyphens | `azure-cost-optimization` | `azure_cost_optimization` |
| No start/end hyphen | `azure-deploy` | `-azure-deploy`, `azure-deploy-` |
| No consecutive hyphens | `azure-deploy` | `azure--deploy` |
| Matches directory | `skill-name` = folder name | Mismatch |
| Length ≤ 64 | 20 chars ✓ | 65+ chars ✗ |

### 2. Description Length

| Score | Length |
|-------|--------|
| Low | < 150 chars |
| Medium | 150-250 chars |
| Medium-High | 250-500 chars |
| Ideal | 300-600 chars |
| Max | 1024 chars |

**Format Rule:** Descriptions MUST use inline double-quoted strings (`"..."`). Do NOT use `>-` folded scalars (incompatible with skills.sh). Do NOT use `|` literal blocks (preserves newlines). Keep descriptions ≤60 words for cross-model reliability.

### 3. Trigger Phrase Detection

**Positive indicators** (case-insensitive):
- `WHEN:` (preferred — scores higher)
- `USE FOR:`
- `USE THIS SKILL`
- `TRIGGERS:`
- `Trigger phrases include`
- `Activate when`

**Scoring:**
- None found → Low
- Implicit (keywords in description) → Medium
- Explicit (WHEN: list) → Medium-High (preferred)
- Explicit (USE FOR: list) → Medium-High (accepted)

### 4. Anti-Trigger Detection

> ⚠️ **Context-dependent risk.** Anti-trigger clauses like "DO NOT USE FOR:" carry different risk levels depending on the skill set size and deployment context.

**Risk assessment by context:**

| Context | Risk Level | Guidance |
|---------|------------|----------|
| Single skill or small set (1-5 skills) with clear domain boundaries | Low | Anti-triggers are low-risk — domain boundaries are obvious |
| Medium skill set (5-15 skills) with some overlap | Moderate | Anti-trigger keywords start competing with other skills' triggers |
| Large skill set (15+ skills) with overlapping domains | **High** | Keyword contamination is measurable — negative keywords become activation keywords on fast-pattern-matching models |

**Why large skill sets are risky:** On Claude Sonnet and similar models that use fast pattern matching (first ~20 words), `DO NOT USE FOR: Function apps` causes Sonnet to key on "Function apps" and **activate** the skill for Functions queries. This was empirically demonstrated across 24 Azure skills ([analysis](https://gist.github.com/kvenkatrajan/52e6e77f5560ca30640490b4cc65d109)). Anthropic's own published skills confirm this pattern — 4 of 5 skills in `anthropics/skills` use positive-only routing.

> **Note:** Anthropic's [Complete Guide to Building Skills](https://resources.anthropic.com/hubfs/The-Complete-Guide-to-Building-Skill-for-Claude.pdf) recommends negative triggers for overtriggering (p25). This is reasonable for small, isolated skill sets. For multi-skill production environments like this repo, Sensei recommends positive routing with `WHEN:` and distinctive quoted phrases as the cross-model-safe alternative.

**Legacy indicators** (still detected, trigger context-dependent warning):
- `DO NOT USE FOR:`
- `NOT FOR:`
- `Don't use this skill`
- `Instead use`
- `Defer to`

**Scoring:**
- Present in large skill sets → emits cross-model compatibility warning
- Present in small skill sets → informational note only
- Absent → no penalty (preferred for cross-model compatibility)

### 5. Compatibility Field

Per the spec, `compatibility` is optional (max 500 characters). Indicates environment requirements.

**What to include:**
- Required tools (azd, az cli, Docker)
- Supported frameworks (.NET 6+, Node.js 18+)
- Required Azure resources
- Optional dependencies

**Example:**
```yaml
compatibility: |
  Requires: Azure CLI, azd CLI
  Supports: Node.js, Python, .NET, Java
  Optional: Docker (for containerized apps)
```

### 6. Optional Spec Fields

The [agentskills.io spec](https://agentskills.io/specification) defines additional optional fields. Sensei **strongly recommends** `license` and `metadata.version` — report a suggestion in the summary if either is missing.

| Field | Spec Status | Sensei Policy |
|-------|-------------|---------------|
| `license` | Optional | **Strongly recommended.** Report suggestion if missing. |
| `metadata.version` | Optional | **Strongly recommended.** Report suggestion if missing. |
| `metadata.*` (other) | Optional | Preserve if present, do not require. |
| `allowed-tools` | Experimental | Preserve if present, do not require. |

> ⚠️ **Warning:** When improving frontmatter, never remove these fields if they already exist.

**Example suggestions in summary:**
```
SUGGESTIONS:
• Add license field (e.g., license: MIT)
• Add metadata.version field (e.g., metadata: { version: "1.0.0" })
```

### 7. SKILL.md Size Limits

Per the spec, SKILL.md should follow progressive disclosure:

| Metric | Limit | Type |
|--------|-------|------|
| SKILL.md tokens | 500 | Soft limit |
| SKILL.md tokens | 5000 | Hard limit |
| SKILL.md lines | 500 | Spec recommendation |
| Description chars | 1024 | Hard limit |
| references/*.md tokens | 1000 | Per-file soft limit |

**Line count check:** The spec recommends keeping SKILL.md under 500 lines. Report a warning if exceeded.

### 8. YAML Description Safety

Descriptions containing YAML special characters (especially `: ` colon-space) **must** use:
- Double-quoted string (`"..."`) — **required format**

Do NOT use `>-` folded scalars (incompatible with skills.sh per [microsoft/GitHub-Copilot-for-Azure#1038](https://github.com/microsoft/GitHub-Copilot-for-Azure/pull/1038)).

| Check | Pass | Fail |
|-------|------|------|
| Uses `"..."` | `description: "Deploy apps..."` | `description: >-` |
| No `: ` in plain value | `description: "Simple text"` | `description: USE FOR: something` |

**Scoring impact:**
- Plain description with `: ` → **Invalid** (will fail to parse)
- Description > 200 chars without `>-` → **Warning** (maintainability concern)

### 9. Security Restrictions

Per Anthropic's [Complete Guide](https://resources.anthropic.com/hubfs/The-Complete-Guide-to-Building-Skill-for-Claude.pdf):

| Check | What it validates |
|-------|-------------------|
| No XML tags | Frontmatter must not contain `<` or `>` (security — frontmatter appears in system prompt) |
| No reserved prefixes | Name must not start with `claude-` or `anthropic-` (reserved by Anthropic) |

---

## Advisory Checks

**Advisory only** — informational, does not change scoring levels.

### 10. Body Structure Quality

Checks whether the SKILL.md body follows Anthropic's recommended structure for effective instructions.

**Sub-checks:**
- **Actionable instructions** — Body uses specific commands, code examples, or step-by-step guidance
- **Examples section** — Has at least one example scenario
- **Error handling** — Documents common failure modes and recovery steps

### 11. Body Progressive Disclosure

Checks whether SKILL.md properly uses progressive disclosure — keeping core instructions in SKILL.md and detailed reference material in `references/`.

**Flag when:**
- SKILL.md body exceeds 500 lines (spec recommends under 500)
- Large code blocks (> 50 lines) that could be in `references/` or `scripts/`

---

## Scoring Algorithm

```
function scoreSkill(skill):
    # Validate name per agentskills.io spec (fail-fast)
    if not isValidName(skill.name):
        report "INVALID: name fails spec validation"
        return "Invalid"
    
    # Check YAML description safety
    if isPlainUnquoted(skill.rawDescription) AND contains(skill.description, ": "):
        report "INVALID: plain description contains ': ' — use double-quoted string"
        return "Invalid"
    
    score = "Low"
    
    # Check description length
    if skill.description.length > 1024:
        report "INVALID: description exceeds 1024-char spec limit"
        return "Invalid"
    if skill.description.length < 150:
        return "Low"
    
    # Check for trigger phrases (WHEN: preferred, USE FOR: accepted)
    hasTriggers = containsTriggerPhrases(skill.description)
    if hasTriggers:
        score = "Medium"
    
    # Check word count for cross-model density
    if wordCount(skill.description) <= 60:
        if hasTriggers:
            score = "Medium-High"
    
    # Warn on anti-triggers (keyword contamination risk)
    if containsAntiTriggers(skill.description):
        warn "DO NOT USE FOR: causes keyword contamination on Sonnet — remove"
    
    # Check for compatibility
    hasCompatibility = skill.compatibility != null
    if score == "Medium-High" AND hasCompatibility:
        score = "High"
    
    return score

function isValidName(name):
    if name.length < 1 OR name.length > 64:
        return false
    if not matches(/^[a-z0-9][a-z0-9-]*[a-z0-9]$|^[a-z0-9]$/):
        return false    # start/end with alphanumeric
    if contains("--"):
        return false    # no consecutive hyphens
    if name != parentDirectoryName:
        return false
    return true

function collectSuggestions(skill):
    suggestions = []
    if skill.license == null:
        suggestions.add("Add license field (e.g., license: MIT)")
    if skill.metadata == null OR skill.metadata.version == null:
        suggestions.add("Add metadata.version field (e.g., metadata: { version: \"1.0.0\" })")
    if usesBlockScalar(skill.rawDescription):
        suggestions.add("Use inline double-quoted string for description (>- incompatible with skills.sh)")
    if containsAntiTriggers(skill.description):
        suggestions.add("Remove DO NOT USE FOR: — causes keyword contamination on Claude Sonnet")
    if wordCount(skill.description) > 60:
        suggestions.add("Trim description to ≤60 words for cross-model reliability")
    return suggestions
```

---

## Current Audit Results

From the [frontmatter audit](https://gist.github.com/spboyer/28c31bf0cafb87489406832633aa31a7):

| Metric | Count | Percentage |
|--------|-------|------------|
| Total Skills | 26 | 100% |
| High Adherence | 0 | 0% |
| Medium Adherence | 14 | 54% |
| Low Adherence | 12 | 46% |

### Low-Adherence Skills (Priority)

1. `appinsights-instrumentation` (71 chars - too brief)
2. `azure-diagnostics`
3. `azure-security`
4. `azure-security-hardening`
5. `azure-observability`
6. `azure-storage`
7. `azure-ai`
8. `azure-validation`
9. `azure-nodejs-production`
10. `entra-app-registration`
11. `azure-rbac`

### Medium-Adherence Skills

1. `azure-deploy` (excellent triggers, missing anti-triggers)
2. `azure-create-app`
3. `azure-deployment-preflight`
4. `azure-postgres`
5. `azure-functions`
6. `azure-quick-review`
7. `azure-cost-optimization`
8. `azure-kusto`
9. `azure-keyvault-expiration-audit`
10. `azure-aigateway`
11. `azure-resource-visualizer`
12. `microsoft-foundry`
13. `skill-authoring`
14. `markdown-token-optimizer`
