## 1. Terminal Output Management

- [x] 1.1 Create `pkg/ink/log_update.go` with terminal output manager
- [x] 1.2 Implement clear previous lines functionality
- [x] 1.3 Implement ANSI cursor movement sequences
- [x] 1.4 Add write output with diff optimization

## 2. Render Options

- [x] 2.1 Create `pkg/ink/options.go` with RenderOptions struct
- [x] 2.2 Implement applyDefaults function for options
- [x] 2.3 Add stdout/stderr writer configuration
- [x] 2.4 Add MaxFps, Interactive, AlternateScreen options

## 3. Ink Instance Core

- [x] 3.1 Create `pkg/ink/ink.go` with Ink struct
- [x] 3.2 Implement NewInk constructor
- [x] 3.3 Add render loop with throttle ticker
- [x] 3.4 Implement onRender method with diff and output
- [x] 3.5 Implement Unmount method with terminal cleanup
- [x] 3.6 Add exit handling with exitChan
- [x] 3.7 Implement WaitUntilExit method

## 4. Render Entry Point

- [x] 4.1 Create `pkg/ink/render.go` with Render function
- [x] 4.2 Create Instance struct with public methods
- [x] 4.3 Implement signal handling (Ctrl+C)
- [x] 4.4 Add Rerender method to Instance
- [x] 4.5 Add Cleanup method to Instance
- [x] 4.6 Add Clear method to Instance
- [x] 4.7 Add WaitUntilRenderFlush method to Instance

## 5. Animation System

- [x] 5.1 Create `pkg/input/animation.go` with AnimationController
- [x] 5.2 Implement Play and Pause methods
- [x] 5.3 Implement Frame method returning current frame
- [x] 5.4 Create UseAnimation hook function
- [x] 5.5 Implement animation subscriber registration
- [x] 5.6 Integrate animation ticks with render loop

## 6. App Instance Hook

- [x] 6.1 Create `pkg/input/app.go` with AppController
- [x] 6.2 Implement UseApp hook function
- [x] 6.3 Add Exit method to AppController
- [x] 6.4 Connect UseApp to Ink instance context

## 7. Examples Update

- [x] 7.1 Update `examples/animation-demo/main.go` with real animation loop
- [x] 7.2 Add spinner example using UseAnimation
- [x] 7.3 Add progress bar example with dynamic updates
- [x] 7.4 Update example to use ink.Render instead of app.Render
- [x] 7.5 Add exit handling with UseInput and UseApp

## 8. Testing and Verification

- [x] 8.1 Verify animation-demo shows real spinner animation
- [x] 8.2 Verify Ctrl+C exits cleanly
- [x] 8.3 Verify cursor is hidden/shown correctly
- [x] 8.4 Verify FPS throttling works as expected
- [x] 8.5 Add unit tests for Ink instance
- [x] 8.6 Add unit tests for AnimationController