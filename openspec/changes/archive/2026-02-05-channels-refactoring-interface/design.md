# Design: Refactor Channels for Testability

## Context
The current channel implementations (`discord`, `slack`, `telegram`) directly use the concrete client structs provided by their respective libraries. This tight coupling makes it impossible to write clean unit tests with mocks, forcing reliance on integration tests or partial struct mocking which is fragile.

## Goals / Non-Goals

**Goals:**
- Decouple `Channel` structs from concrete external library types.
- Define internal interfaces that cover only the used subset of methods.
- Implement thin adapters that wrap the external libraries.
- Enable full unit testing of `Channel` logic using generated or manual mocks.

**Non-Goals:**
- Changing the external libraries themselves.
- Changing the public API of the `Channel` packages (other than internal structure).
- Adding new features to the channels.

## Decisions

### 1. Interface-Adapter Pattern
We will define interfaces within each package (e.g., `discord.Session`) that mirror the external library methods we use.
We will then create an adapter struct (e.g., `discordSession`) that embeds or holds the real client and implements the interface.
The `Channel` struct will hold the interface, allowing `New()` to inject the adapter and tests to inject a mock.

### 2. Interface Scope
Interfaces will only define the methods currently in use. We will not attempt to wrap the entire external library API, as that would be a large maintenance burden.

## Risks / Trade-offs

- **Maintenance**: If the external library changes its API, we must update both our adapter and our interface definition. This adds a layer of maintenance but is standard for the adapter pattern.
- **Performance**: Negligible overhead from interface method calls.
