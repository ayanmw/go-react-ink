# render-loop Specification

## Purpose
TBD - created by archiving change complete-tui-animation-effects. Update Purpose after archive.
## Requirements
### Requirement: Render loop calculates layout
The system SHALL calculate layout on each render tick using the fiber reconciler.

#### Scenario: Layout calculation
- **WHEN** render loop tick occurs
- **THEN** system reconciles the element tree and calculates layout

### Requirement: Render loop diffs output
The system SHALL diff the current output with the previous output and only write changes.

#### Scenario: Output unchanged
- **WHEN** rendered output is identical to previous frame
- **THEN** system skips writing to stdout

#### Scenario: Output changed
- **WHEN** rendered output differs from previous frame
- **THEN** system writes the diff using ANSI escape sequences

### Requirement: Render loop handles animation ticks
The system SHALL check for animation subscribers on each tick and advance animations.

#### Scenario: Animation tick with subscribers
- **WHEN** there are active animation subscribers
- **THEN** system advances animation frame and triggers re-render

#### Scenario: Animation tick without subscribers
- **WHEN** there are no active animation subscribers
- **THEN** system skips animation tick processing

### Requirement: Render loop handles exit signal
The system SHALL stop the render loop when exit is signaled.

#### Scenario: Exit signal received
- **WHEN** exitChan receives a value
- **THEN** system finishes unmount and exits the render loop

### Requirement: Render loop uses ticker for throttling
The system SHALL use a time.Ticker to enforce FPS throttling.

#### Scenario: Throttle ticker fires
- **WHEN** throttle ticker fires
- **THEN** system processes any pending renders

