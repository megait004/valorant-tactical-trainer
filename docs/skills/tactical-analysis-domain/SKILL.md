---
name: tactical-analysis-domain
description: Use when changing match analysis, tactical scoring, findings, recommendations, training plans, or report generation.
---

# Tactical Analysis Domain

## Use When

- Adding or changing tactical metrics.
- Adding analysis findings.
- Changing recommendation rules.
- Working with match, round, economy, plant, defuse, map, or agent analysis.
- Generating training plans.

## Rules

- MVP analysis must be deterministic and rule-based.
- Every finding must include evidence.
- Every recommendation must link to a finding.
- Every recommendation must include reason, severity/confidence context, and a drill.
- Domain logic must not depend on React, Wails, SQLite, or provider DTOs.
- Missing data should lower confidence or skip a rule, not crash analysis.
- Prefer small pure functions for metric calculations.

## Evidence Shape

Evidence may include:

- match ids.
- round numbers.
- map.
- agent.
- metric value.
- sample size.
- comparison baseline.

## Validation

- Use table-driven Go tests for metric and rule behavior.
- Include tests for missing or partial data.
- Keep rule thresholds explicit and documented.
