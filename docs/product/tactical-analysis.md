# Tactical Analysis

## Goal

The analysis engine converts stored match data into explainable tactical findings and training recommendations.

MVP analysis is rule-based and deterministic.

## Domain Rules

- Every finding must include evidence.
- Every recommendation must link to a finding.
- Findings must include severity and confidence.
- Missing provider fields must degrade confidence instead of crashing analysis.
- Domain logic must not depend on React, Wails, SQLite, or provider DTOs.

## Metrics

Player metrics:

- KDA.
- Kill/death ratio.
- Assist rate.
- Headshot percentage.
- Damage made and received.
- Agent performance.
- Map performance.
- Recent trend.

Round metrics:

- Round win/loss.
- Plant outcomes.
- Defuse outcomes.
- Plant site distribution.
- Economy effectiveness.
- Ability casts per round.
- AFK rounds, friendly fire, and rounds in spawn when available.

Map metrics:

- Win rate by map.
- Weakest map.
- Strongest map.
- Site pattern when plant data exists.
- Agent-map pairing when sample size is enough.

Economy metrics:

- Average spend.
- Average loadout value.
- Low-spend conversion.
- High-spend round losses.

## Finding Model

```text
type
severity
confidence
title
description
evidence
suggested_action
```

Severity values:

- `low`
- `medium`
- `high`
- `critical`

Confidence is a number from `0.0` to `1.0`.

Evidence may include:

- match ids.
- round numbers.
- map.
- agent.
- metric value.
- sample size.
- comparison baseline.

## Recommendation Model

```text
title
drill
priority
reason
linked_finding
evidence
status
```

Status values:

- `new`
- `planned`
- `in_progress`
- `done`
- `ignored`

## MVP Rules

Low headshot percentage:

```text
If headshot_percentage < threshold and sample_size >= minimum:
  create aim discipline recommendation
```

Weak map:

```text
If map winrate is low and match_count >= minimum:
  create map-specific review recommendation
```

Poor plant conversion:

```text
If bomb planted but round lost often:
  create post-plant positioning/utility recommendation
```

Poor defuse success:

```text
If opponent plant rounds frequently lost:
  create retake drill recommendation
```

Low ability usage:

```text
If ability casts per round are below role baseline:
  create utility usage recommendation
```

Economy issue:

```text
If high-spend rounds are lost frequently:
  create buy-round trading/spacing recommendation
```

## Virtual Tactical Assistant MVP

The assistant provides tactical lookup cards rather than live game automation. Users choose map, agent, side, and phase manually. The app combines local cards with simple economy rules to suggest safe actions.

Supported MVP card categories:

- `composition`
- `default-strat`
- `crosshair`
- `lineup`

Economy rules:

- credits below 2000 -> Eco.
- credits below 3900 after a loss -> Light / Half Buy.
- credits at or above 3900 -> Full Buy.
- otherwise -> Force Buy only if the whole team commits.

The assistant must not depend on memory reading, process injection, or hidden match-state detection.
