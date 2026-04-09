## ADDED Requirements

### Requirement: Render function creates application instance
The system SHALL provide a `Render(element, options)` function that creates an application instance and starts the render loop.

#### Scenario: Basic render call
- **WHEN** user calls `ink.Render(<App />, nil)`
- **THEN** system returns an Instance with Unmount, WaitUntilExit, Rerender, and Cleanup methods

#### Scenario: Render with options
- **WHEN** user calls `ink.Render(<App />, &ink.RenderOptions{Stdout: customWriter})`
- **THEN** system uses the provided options for rendering

### Requirement: Render starts render loop in goroutine
The system SHALL start the render loop in a separate goroutine when Render is called.

#### Scenario: Render loop starts automatically
- **WHEN** user calls `ink.Render(<App />, nil)`
- **THEN** system starts a goroutine that continuously renders the application

### Requirement: Render supports max FPS throttling
The system SHALL throttle rendering to a configurable maximum FPS (default: 30).

#### Scenario: Default FPS throttle
- **WHEN** user renders without specifying MaxFps
- **THEN** system throttles rendering to 30 FPS maximum

#### Scenario: Custom FPS throttle
- **WHEN** user renders with `RenderOptions{MaxFps: 60}`
- **THEN** system throttles rendering to 60 FPS maximum

### Requirement: Instance provides unmount method
The system SHALL provide an Unmount method that cleanly stops the application.

#### Scenario: Clean unmount
- **WHEN** user calls `instance.Unmount()`
- **THEN** system stops the render loop and restores terminal state

#### Scenario: Unmount with result
- **WHEN** user calls `instance.Unmount(result)`
- **THEN** system stores the result and returns it from WaitUntilExit

### Requirement: Instance provides wait until exit
The system SHALL provide a WaitUntilExit method that blocks until the application exits.

#### Scenario: Wait for exit
- **WHEN** user calls `instance.WaitUntilExit()`
- **THEN** system blocks until Unmount is called and returns the exit result

### Requirement: Instance provides rerender method
The system SHALL provide a Rerender method that triggers a re-render with a new element tree.

#### Scenario: Rerender with new content
- **WHEN** user calls `instance.Rerender(<NewApp />)`
- **THEN** system updates the rendered output to reflect the new element tree
