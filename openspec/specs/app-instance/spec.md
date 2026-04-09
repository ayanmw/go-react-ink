# app-instance Specification

## Purpose
TBD - created by archiving change complete-tui-animation-effects. Update Purpose after archive.
## Requirements
### Requirement: UseApp returns app controller
The system SHALL provide a UseApp hook that returns an AppController.

#### Scenario: Basic app hook
- **WHEN** component calls `input.UseApp(ctx)`
- **THEN** system returns an AppController with Exit method

### Requirement: App controller provides Exit method
The system SHALL provide an Exit method that triggers application exit.

#### Scenario: Exit from component
- **WHEN** user calls `app.Exit()`
- **THEN** system calls Unmount on the application instance

#### Scenario: Exit with result
- **WHEN** user calls `app.Exit(result)`
- **THEN** system calls Unmount with the result

### Requirement: Instance provides Cleanup method
The system SHALL provide a Cleanup method for manual resource cleanup.

#### Scenario: Manual cleanup
- **WHEN** user calls `instance.Cleanup()`
- **THEN** system releases all resources without waiting for exit

### Requirement: Instance provides Clear method
The system SHALL provide a Clear method to clear the terminal output.

#### Scenario: Clear terminal
- **WHEN** user calls `instance.Clear()`
- **THEN** system clears all rendered content from terminal

### Requirement: Instance provides WaitUntilRenderFlush
The system SHALL provide a WaitUntilRenderFlush method that waits for pending renders.

#### Scenario: Wait for render flush
- **WHEN** user calls `instance.WaitUntilRenderFlush()`
- **THEN** system blocks until all pending renders are complete

