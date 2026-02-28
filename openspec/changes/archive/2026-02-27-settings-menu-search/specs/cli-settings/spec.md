## ADDED Requirements

### Requirement: Grouped Section Layout
The settings menu SHALL organize categories into named sections. Each section SHALL have a title header rendered above its categories with a visual separator line between sections.

The sections SHALL be, in order:
1. **Core** — Providers, Agent, Server, Session
2. **Communication** — Channels, Tools, Multi-Agent, A2A Protocol
3. **AI & Knowledge** — Knowledge, Skill, Observational Memory, Embedding & RAG, Graph Store, Librarian
4. **Infrastructure** — Payment, Cron Scheduler, Background Tasks, Workflow Engine
5. **P2P Network** — P2P Network, P2P ZKP, P2P Pricing, P2P Owner Protection, P2P Sandbox
6. **Security** — Security, Auth, Security Keyring, Security DB Encryption, Security KMS
7. *(untitled)* — Save & Exit, Cancel

#### Scenario: Section headers displayed
- **WHEN** user views the settings menu in normal (non-search) mode
- **THEN** named section headers SHALL be rendered above each group of categories with separator lines between sections

#### Scenario: Flat cursor across sections
- **WHEN** user navigates with arrow keys
- **THEN** the cursor SHALL move through all categories across sections as a flat list, skipping section headers

### Requirement: Keyword Search
The settings menu SHALL support real-time keyword search to filter categories.

#### Scenario: Activate search
- **WHEN** user presses `/` in normal mode
- **THEN** the menu SHALL enter search mode, display a focused text input with `/ ` prompt and "Type to search..." placeholder, and reset the cursor to 0

#### Scenario: Filter categories
- **WHEN** user types a search query
- **THEN** the menu SHALL filter categories by case-insensitive substring match against title, description, and ID, updating results in real-time

#### Scenario: Empty search query
- **WHEN** the search input is empty or whitespace-only
- **THEN** all categories SHALL be displayed (no filtering)

#### Scenario: No results
- **WHEN** the search query matches no categories
- **THEN** the menu SHALL display "No matching items" in muted italic text

#### Scenario: Select from search results
- **WHEN** user presses Enter during search mode
- **THEN** the selected filtered category SHALL be activated, search mode SHALL exit, and the search input SHALL be cleared

#### Scenario: Cancel search
- **WHEN** user presses Esc during search mode
- **THEN** search mode SHALL be cancelled, the filtered list SHALL be cleared, and the full grouped menu SHALL be restored

#### Scenario: Navigate search results
- **WHEN** user presses up/down (or shift+tab/tab) during search mode
- **THEN** the cursor SHALL move within the filtered results list

### Requirement: Search Match Highlighting
The settings menu SHALL highlight matching substrings in search results.

#### Scenario: Highlight matching text
- **WHEN** categories are displayed during an active search with a non-empty query
- **THEN** the first matching substring in each category's title and description SHALL be rendered in amber/warning color with bold styling

#### Scenario: Selected item highlight
- **WHEN** the cursor is on a filtered category during search
- **THEN** the matching substring SHALL additionally be underlined

### Requirement: Search Help Bar
The help bar SHALL update based on the current mode.

#### Scenario: Normal mode help bar
- **WHEN** the menu is in normal mode
- **THEN** the help bar SHALL display: Navigate, Select, Search (`/`), Back (`Esc`)

#### Scenario: Search mode help bar
- **WHEN** the menu is in search mode
- **THEN** the help bar SHALL display: Navigate, Select, Cancel (`Esc`)
