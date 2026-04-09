# animation-system Specification

## Purpose
TBD - created by archiving change complete-tui-animation-effects. Update Purpose after archive.
## Requirements
### Requirement: UseAnimation returns controller
The system SHALL provide a UseAnimation hook that returns an AnimationController.

#### Scenario: Basic animation hook
- **WHEN** component calls `input.UseAnimation(ctx, 10)`
- **THEN** system returns an AnimationController with Play, Pause, and Frame methods

### Requirement: Animation controller provides Play method
The system SHALL provide a Play method that starts the animation.

#### Scenario: Start animation
- **WHEN** user calls `anim.Play()`
- **THEN** system sets isPlaying to true and records start time

### Requirement: Animation controller provides Pause method
The system SHALL provide a Pause method that stops the animation.

#### Scenario: Pause animation
- **WHEN** user calls `anim.Pause()`
- **THEN** system sets isPlaying to false

### Requirement: Animation controller provides Frame method
The system SHALL provide a Frame method that returns the current frame number.

#### Scenario: Get current frame
- **WHEN** user calls `anim.Frame()`
- **THEN** system returns the current frame counter value

### Requirement: Animation integrates with render loop
The system SHALL register animation controllers with the global animation manager.

#### Scenario: Animation triggers re-render
- **WHEN** animation frame advances
- **THEN** system schedules a re-render

### Requirement: Animation respects FPS setting
The system SHALL advance animation frames at the specified FPS rate.

#### Scenario: 10 FPS animation
- **WHEN** user creates animation with `UseAnimation(ctx, 10)`
- **THEN** system advances frame approximately 10 times per second

