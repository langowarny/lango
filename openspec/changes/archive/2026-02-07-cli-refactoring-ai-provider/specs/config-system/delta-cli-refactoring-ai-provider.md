# Delta Spec: Config System - Multi-Provider CLI

## Overview

Updates the config system to support CLI-generated multi-provider configurations.

## MODIFIED Requirements

### Requirement: Provider Configuration Precedence
**Reason**: To align CLI behavior with Core behavior.
**Impact**: Configuration validation SHALL NOT fail if legacy `agent.provider` fields are missing, provided a valid provider is defined in the `providers` map and referenced by `agent.provider`.
