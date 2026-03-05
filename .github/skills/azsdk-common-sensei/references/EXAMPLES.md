# Examples

Before and after examples of frontmatter improvements, including token counts.

## Example 1: appinsights-instrumentation (Low → Medium-High)

### Before (Low Adherence)

**SKILL.md frontmatter:**
```yaml
---
name: appinsights-instrumentation
description: 'Instrument a webapp to send useful telemetry data to Azure App Insights'
---
```

**Metrics:**
- Score: Low
- Tokens: ~142
- Triggers: 0
- Anti-triggers: 0

**Problems:**
- ❌ Only 71 characters - too brief
- ❌ No trigger phrases
- ❌ No anti-triggers
- ❌ Agent doesn't know when to activate

**triggers.test.ts:**
```javascript
const shouldTriggerPrompts = [
  // Empty - no test coverage
];

const shouldNotTriggerPrompts = [
  'What is the weather today?',
  'Help me write a poem',
];
```

### After (Medium-High Adherence)

**SKILL.md frontmatter:**
```yaml
---
name: appinsights-instrumentation
description: "Instrument web applications to send telemetry to Azure Application Insights for monitoring and diagnostics. WHEN: \"add App Insights\", \"instrument my app\", \"set up application monitoring\", \"add telemetry\", \"track requests and dependencies\", \"ASP.NET Core telemetry\", \"Node.js Application Insights\"."
---
```

**Metrics:**
- Score: Medium-High ✅
- Tokens: ~385 (under 500 budget ✅)
- Triggers: 7 (WHEN: format)

**Improvements:**
- ✅ Informative but concise (under 60 words)
- ✅ Uses inline double-quoted string (skills.sh compatible)
- ✅ Explicit "WHEN:" trigger phrases with quoted keywords
- ✅ No "DO NOT USE FOR:" (avoids keyword contamination)

**triggers.test.ts:**
```javascript
const shouldTriggerPrompts = [
  'Add Application Insights to my web app',
  'Instrument my ASP.NET Core application for monitoring',
  'Set up telemetry for my Node.js app',
  'How do I track requests in App Insights?',
  'Add Application Insights monitoring to my project',
  'Configure App Insights for my Azure web app',
];

const shouldNotTriggerPrompts = [
  'What is the weather today?',
  'Help me write a poem',
  'Query my Application Insights logs',  // → azure-observability
  'Create an alert in Azure Monitor',     // → azure-observability
  'Show me my App Insights dashboard',    // → azure-observability
  'How much does App Insights cost?',     // → azure-cost-optimization
  'Help me with AWS CloudWatch',          // Wrong cloud provider
];
```

---

## Example 2: azure-security (Low → Medium-High)

### Before (Low Adherence)

**SKILL.md frontmatter:**
```yaml
---
name: azure-security
description: 'Azure Security Services including Key Vault, Managed Identity, RBAC, Entra ID, and Defender. Provides secrets management, credential-free authentication, role-based access control, and threat protection.'
---
```

**Problems:**
- ❌ Catalog-style description (lists services, not actions)
- ❌ No trigger phrases
- ❌ No anti-triggers
- ❌ Overlaps with azure-security-hardening, entra-app-registration

### After (Medium-High Adherence)

**SKILL.md frontmatter:**
```yaml
---
name: azure-security
description: "Overview of Azure security services and concepts including Key Vault, Managed Identity, RBAC, Entra ID, and Defender for Cloud. WHEN: \"Azure security overview\", \"what security services are available\", \"explain managed identity\", \"RBAC basics\", \"Key Vault concepts\", \"Entra ID overview\", \"Defender for Cloud features\"."
---
```

**Improvements:**
- ✅ Reframed as "overview/concepts" skill
- ✅ Explicit WHEN: triggers for educational queries
- ✅ No "DO NOT USE FOR:" (avoids keyword contamination with azure-security-hardening, entra-app-registration)

---

## Example 3: azure-deploy (Medium → Medium-High)

### Before (Medium Adherence)

**SKILL.md frontmatter:**
```yaml
---
name: azure-deploy
description: >-
  Deploy applications to Azure App Service, Azure Functions, and Static Web Apps.
  USE THIS SKILL when users want to deploy, publish, host, or run their application
  on Azure. This skill detects application type and recommends optimal Azure service.
  Trigger phrases include "deploy to Azure", "host on Azure", "publish to Azure".
---
```

**Status:**
- ✅ Good trigger phrases
- ❌ Missing anti-triggers (collision with azure-create-app, azure-deployment-preflight)

### After (Medium-High Adherence)

**SKILL.md frontmatter:**
```yaml
---
name: azure-deploy
description: "Deploy applications to Azure App Service, Azure Functions, and Static Web Apps. WHEN: \"deploy to Azure\", \"host on Azure\", \"publish to Azure\", \"run on Azure\", \"azd up\", \"deploy my app\", \"push to Azure\", \"get this running in the cloud\"."
---
```

**Improvements:**
- ✅ Uses WHEN: with distinctive quoted trigger phrases
- ✅ No "DO NOT USE FOR:" — previously mentioned "azure-create-app" and "azure-deployment-preflight" which caused keyword contamination
- ✅ Concise and cross-model optimized

---

## Example 4: Test Prompt Patterns

### Good Trigger Prompts

```javascript
// Specific and actionable
'Deploy my Node.js app to Azure App Service'
'How do I publish a React app to Azure Static Web Apps?'
'Set up Azure Functions for my Python project'

// Uses skill keywords
'Run azd up to deploy my application'
'Host my web app on Azure'

// Natural variations
'I need to get my app running on Azure'
'Can you help me deploy to the cloud?'
```

### Good Anti-Trigger Prompts

```javascript
// Unrelated topics (always include some)
'What is the weather today?',
'Help me write a poem about clouds',
'Explain quantum computing',

// Other cloud providers
'Deploy to AWS Lambda',
'Host on Google Cloud Run',
'Push to Heroku',

// Related but different Azure services
'Validate my Bicep template',        // → azure-deployment-preflight
'Create an azure.yaml file',         // → azure-create-app
'Why did my deployment fail?',       // → azure-diagnostics
'How much will this deployment cost?', // → azure-cost-optimization
```

---

## Example 5: Commit Message

**Good commit message:**
```
sensei: improve appinsights-instrumentation frontmatter

- Added USE FOR trigger phrases: "add App Insights", "instrument my app", etc.
- Added DO NOT USE FOR anti-triggers: querying logs, creating alerts
- Updated triggers.test.ts with 6 should-trigger and 7 should-not-trigger prompts
- Score improved: Low → Medium-High
- All tests passing
```

**Minimal commit message (also acceptable):**
```
sensei: improve appinsights-instrumentation frontmatter
```

---

## Anti-Pattern Examples

### ❌ Don't Do This: Anti-Trigger Keyword Contamination

```yaml
description: "Deploy to Azure. USE FOR: \"deploy app\", \"push to production\". DO NOT USE FOR: Function apps (use azure-functions), storage accounts (use azure-storage), Kubernetes (use azure-aks)."
```

**Problem:** The "DO NOT USE FOR" clause introduces "Function apps", "storage", "Kubernetes" — on Claude Sonnet these keywords cause this skill to activate for Functions, Storage, and AKS tasks. Remove anti-triggers and use positive routing with distinctive `WHEN:` phrases instead.

### ❌ Don't Do This: Description Too Dense (>60 words)

```yaml
description: "This skill orchestrates the complete Azure deployment workflow including preparation, validation, and execution phases for web applications, APIs, Function apps, container apps, and static web apps. It generates Bicep templates, Terraform configurations, azure.yaml files, and Dockerfiles. USE FOR: deploy, create app, build, migrate, modernize, prepare, validate, provision, configure, set up, initialize..."
```

**Problem:** At 60+ words, Sonnet's attention is diluted across all 24+ skill descriptions. Cap at ≤60 words and front-load the unique signal.

### ❌ Don't Do This: Using >- Folded Scalars

```yaml
description: >-
  Deploy applications to Azure App Service, Azure Functions, and Static Web Apps.
  USE FOR: "deploy to Azure", "host on Azure", "publish to Azure".
```

**Problem:** The `>-` folded scalar syntax is incompatible with skills.sh and other registries. Use inline double-quoted strings instead.

### ❌ Don't Do This: Overly Long Description

```yaml
description: >-
  This skill helps you instrument your web applications to send telemetry data
  to Azure Application Insights which is a feature of Azure Monitor that provides
  extensible application performance management (APM) and monitoring for live
  web applications. It supports a wide variety of platforms including .NET, Node.js,
  Java, and Python. Application Insights can be used to monitor request rates,
  response times, failure rates, dependency tracking, page views, user counts,
  exceptions, and much more. You can use this skill when you want to add App
  Insights to your app, instrument your application, set up monitoring, add
  telemetry, or track requests. Do not use this skill if you want to query logs
  or create dashboards or set up alerts because those are handled by other skills.
  [... continues for 1200+ characters ...]
```

**Problem:** Exceeds 1024 character limit, too verbose.

### ❌ Don't Do This: Vague Anti-Triggers

```yaml
description: >-
  Deploy apps to Azure. USE FOR: deployment.
  DO NOT USE FOR: other stuff.
```

**Problem:** Too vague - agent can't determine what "other stuff" means.

### ❌ Don't Do This: Mismatched Tests

```yaml
# SKILL.md says:
USE FOR: "deploy to Azure", "host on Azure"
```

```javascript
// But triggers.test.ts has:
const shouldTriggerPrompts = [
  'Create a new Azure project',  // Wrong! This is azure-create-app
  'Validate my Bicep file',      // Wrong! This is azure-deployment-preflight
];
```

**Problem:** Test prompts don't match the actual trigger phrases.

---

## Example 6: Sensei Summary Output

After completing improvements, Sensei displays a summary:

```
╔══════════════════════════════════════════════════════════════════╗
║  SENSEI SUMMARY: appinsights-instrumentation                     ║
╠══════════════════════════════════════════════════════════════════╣
║  BEFORE                          AFTER                           ║
║  ──────                          ─────                           ║
║  Score: Low                      Score: Medium-High              ║
║  Tokens: 142                     Tokens: 385                     ║
║  Triggers: 0                     Triggers: 7                     ║
║  Anti-triggers: 0                Anti-triggers: 4                ║
║                                                                  ║
║  TOKEN STATUS: ✅ Under budget (385 < 500)                       ║
║                                                                  ║
║  SUGGESTIONS NOT IMPLEMENTED:                                    ║
║  • (none - skill is under budget and clear)                      ║
║                                                                  ║
║  Choose an action:                                               ║
║    [C] Commit changes                                            ║
║    [I] Create GitHub issue with suggestions                      ║
║    [S] Skip (discard changes)                                    ║
╚══════════════════════════════════════════════════════════════════╝
```

### Summary with Token Suggestions

When token optimizations are available:

```
╔══════════════════════════════════════════════════════════════════╗
║  SENSEI SUMMARY: azure-deploy                                    ║
╠══════════════════════════════════════════════════════════════════╣
║  BEFORE                          AFTER                           ║
║  ──────                          ─────                           ║
║  Score: Medium                   Score: Medium-High              ║
║  Tokens: 623                     Tokens: 589                     ║
║  Triggers: 4                     Triggers: 8                     ║
║  Anti-triggers: 0                Anti-triggers: 3                ║
║                                                                  ║
║  TOKEN STATUS: ⚠️  Above soft limit (589 > 500)                  ║
║                                                                  ║
║  SUGGESTIONS NOT IMPLEMENTED:                                    ║
║  • Remove bullet point decorators (-8 tokens)                    ║
║  • Shorten "In order to deploy" → "To deploy" (-12 tokens)       ║
║  • Consolidate example list into paragraph (-18 tokens)          ║
║  Potential savings: ~38 tokens → 551 tokens                      ║
║                                                                  ║
║  Choose an action:                                               ║
║    [C] Commit changes                                            ║
║    [I] Create GitHub issue with suggestions                      ║
║    [S] Skip (discard changes)                                    ║
╚══════════════════════════════════════════════════════════════════╝
```

**If user chooses [I] Create issue:**

Creates GitHub issue:
- **Title:** `[sensei] Token optimization suggestions for azure-deploy`
- **Labels:** `enhancement`, `skill-quality`
- **Body:**
  ```markdown
  ## Summary
  
  Sensei improved `azure-deploy` frontmatter but identified token optimization opportunities.
  
  | Metric | Before | After |
  |--------|--------|-------|
  | Score | Medium | Medium-High |
  | Tokens | 623 | 589 |
  | Status | ⚠️ Above soft limit | |
  
  ## Suggestions
  
  - [ ] Remove bullet point decorators (-8 tokens)
  - [ ] Shorten "In order to deploy" → "To deploy" (-12 tokens)
  - [ ] Consolidate example list into paragraph (-18 tokens)
  
  **Potential savings:** ~38 tokens → 551 tokens (under 500 soft limit)
  
  ## References
  
  - [OPTIMIZATION-PATTERNS.md](/.github/skills/markdown-token-optimizer/references/OPTIMIZATION-PATTERNS.md)
  ```


---

## Skill Body Patterns

Templates for structuring the SKILL.md body content, adapted from Anthropic's [Complete Guide to Building Skills](https://resources.anthropic.com/hubfs/The-Complete-Guide-to-Building-Skill-for-Claude.pdf).

### Pattern 1: Sequential Workflow
Use when users need multi-step processes in a specific order. Key techniques: explicit step ordering, dependencies between steps, validation at each stage.

### Pattern 2: Multi-MCP Coordination
Use when workflows span multiple services or MCP servers. Key techniques: clear phase separation, data passing between MCPs, validation before next phase.

### Pattern 3: Iterative Refinement
Use when output quality improves with iteration. Generate draft, run validation, address issues, re-validate, repeat until quality threshold met.

### Pattern 4: Context-Aware Tool Selection
Use when the same outcome requires different tools depending on context. Build a decision tree with clear criteria, fallback options, and transparency about choices.

### Pattern 5: Domain-Specific Intelligence
Use when the skill adds specialized knowledge beyond tool access. Apply domain rules before processing, document compliance decisions, maintain audit trail.
