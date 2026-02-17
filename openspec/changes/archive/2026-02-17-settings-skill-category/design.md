## Context

The Skill system was migrated from Ent/SQLite to file-based storage (`skill-file-migration`), introducing `SkillConfig` with `Enabled` and `SkillsDir` fields. The onboard/settings split (`onboard-settings-split`) created a dedicated settings editor with 15+ categories. However, the Skill config was never wired into the settings UI — it remained configurable only through JSON import/export.

## Goals / Non-Goals

**Goals:**
- Add a standalone "Skill" menu category in `lango settings` with Enabled and SkillsDir fields
- Wire `skill_enabled` and `skill_dir` field keys to `config.Skill` in `UpdateConfigFromForm`
- Remove "Skills" from the Knowledge menu description since skills are now independent
- Update README to reflect skill-file-migration and onboard-settings-split changes

**Non-Goals:**
- No changes to the onboard wizard (5-step guided flow covers essentials only)
- No changes to skill system internals or SkillConfig schema
- No skill browsing/management UI (just config fields)

## Decisions

**1. Separate Skill category vs. embedding in Knowledge menu**

Skill is an independent subsystem with its own `SkillConfig` struct, separate from `KnowledgeConfig`. A dedicated menu category keeps the settings editor aligned with the config structure and avoids confusion about the skill/knowledge relationship.

**2. Minimal form fields (Enabled + SkillsDir)**

`SkillConfig` has exactly two fields. No need for additional UI complexity — the form follows the same pattern as other simple categories like Multi-Agent.

## Risks / Trade-offs

- [Menu length] Adding another category increases the menu from 17 to 18 items → Acceptable; users can scroll freely and the Skill entry sits logically after Knowledge.
