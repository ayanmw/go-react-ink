// Package main provides LSP server for .gox files
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/ayanmw/go-react-ink/internal/compiler"
)

func main() {
	server := NewLSPServer()
	if err := server.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "LSP server error: %v\n", err)
		os.Exit(1)
	}
}

// LSPServer implements a basic LSP server
type LSPServer struct {
	reader    *bufio.Reader
	writer    io.Writer
	mu        sync.Mutex
	documents map[string]string // URI -> content
}

// NewLSPServer creates a new LSP server
func NewLSPServer() *LSPServer {
	return &LSPServer{
		reader:    bufio.NewReader(os.Stdin),
		writer:    os.Stdout,
		documents: make(map[string]string),
	}
}

// Run starts the LSP server
func (s *LSPServer) Run() error {
	for {
		msg, err := s.readMessage()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		go s.handleMessage(msg)
	}
}

// Message represents an LSP message
type Message struct {
	Header  map[string]string
	Content []byte
}

func (s *LSPServer) readMessage() (*Message, error) {
	msg := &Message{Header: make(map[string]string)}

	// Read headers
	for {
		line, err := s.reader.ReadString('\r')
		if err != nil {
			return nil, err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) == 2 {
			msg.Header[parts[0]] = parts[1]
		}
	}

	// Read empty line after headers
	_, err := s.reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	// Read content
	var length int
	fmt.Sscanf(msg.Header["Content-Length"], "%d", &length)
	msg.Content = make([]byte, length)
	_, err = io.ReadFull(s.reader, msg.Content)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (s *LSPServer) writeMessage(content []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(content))
	s.writer.Write([]byte(header))
	s.writer.Write(content)
}

func (s *LSPServer) handleMessage(msg *Message) {
	var request struct {
		JSONRPC string          `json:"jsonrpc"`
		ID      interface{}     `json:"id"`
		Method  string          `json:"method"`
		Params  json.RawMessage `json:"params"`
	}

	if err := json.Unmarshal(msg.Content, &request); err != nil {
		return
	}

	switch request.Method {
	case "initialize":
		s.handleInitialize(request.ID)
	case "initialized":
		// No response needed
	case "shutdown":
		s.handleShutdown(request.ID)
	case "exit":
		os.Exit(0)
	case "textDocument/didOpen":
		s.handleDidOpen(request.Params)
	case "textDocument/didChange":
		s.handleDidChange(request.Params)
	case "textDocument/didClose":
		s.handleDidClose(request.Params)
	case "textDocument/diagnostic":
		s.handleDiagnostic(request.ID, request.Params)
	case "textDocument/completion":
		s.handleCompletion(request.ID, request.Params)
	case "textDocument/hover":
		s.handleHover(request.ID, request.Params)
	}
}

func (s *LSPServer) handleInitialize(id interface{}) {
	result := map[string]interface{}{
		"capabilities": map[string]interface{}{
			"textDocumentSync": 1,
			"completionProvider": map[string]interface{}{
				"triggerCharacters": []string{"<", "{"},
			},
			"hoverProvider": true,
			"diagnosticProvider": map[string]interface{}{
				"interFileDependencies": false,
				"workspaceDiagnostics":  false,
			},
		},
		"serverInfo": map[string]interface{}{
			"name":    "gox-lsp",
			"version": "0.1.0",
		},
	}

	s.sendResponse(id, result)
}

func (s *LSPServer) handleShutdown(id interface{}) {
	s.sendResponse(id, nil)
}

func (s *LSPServer) handleDidOpen(params json.RawMessage) {
	var notification struct {
		TextDocument struct {
			URI     string `json:"uri"`
			Text    string `json:"text"`
			Version int    `json:"version"`
		} `json:"textDocument"`
	}

	if err := json.Unmarshal(params, &notification); err != nil {
		return
	}

	s.mu.Lock()
	s.documents[notification.TextDocument.URI] = notification.TextDocument.Text
	s.mu.Unlock()

	s.publishDiagnostics(notification.TextDocument.URI, notification.TextDocument.Text)
}

func (s *LSPServer) handleDidChange(params json.RawMessage) {
	var notification struct {
		TextDocument struct {
			URI     string `json:"uri"`
			Version int    `json:"version"`
		} `json:"textDocument"`
		ContentChanges []struct {
			Text string `json:"text"`
		} `json:"contentChanges"`
	}

	if err := json.Unmarshal(params, &notification); err != nil {
		return
	}

	if len(notification.ContentChanges) > 0 {
		text := notification.ContentChanges[0].Text
		s.mu.Lock()
		s.documents[notification.TextDocument.URI] = text
		s.mu.Unlock()

		s.publishDiagnostics(notification.TextDocument.URI, text)
	}
}

func (s *LSPServer) handleDidClose(params json.RawMessage) {
	var notification struct {
		TextDocument struct {
			URI string `json:"uri"`
		} `json:"textDocument"`
	}

	if err := json.Unmarshal(params, &notification); err != nil {
		return
	}

	s.mu.Lock()
	delete(s.documents, notification.TextDocument.URI)
	s.mu.Unlock()
}

func (s *LSPServer) handleDiagnostic(id interface{}, params json.RawMessage) {
	var request struct {
		TextDocument struct {
			URI string `json:"uri"`
		} `json:"textDocument"`
	}

	if err := json.Unmarshal(params, &request); err != nil {
		s.sendResponse(id, nil)
		return
	}

	s.mu.Lock()
	text := s.documents[request.TextDocument.URI]
	s.mu.Unlock()

	diagnostics := s.getDiagnostics(text)
	s.sendResponse(id, map[string]interface{}{
		"kind":     "full",
		"items":    diagnostics,
		"resultId": "",
	})
}

func (s *LSPServer) handleCompletion(id interface{}, params json.RawMessage) {
	var request struct {
		TextDocument struct {
			URI string `json:"uri"`
		} `json:"textDocument"`
		Position struct {
			Line      int `json:"line"`
			Character int `json:"character"`
		} `json:"position"`
	}

	if err := json.Unmarshal(params, &request); err != nil {
		s.sendResponse(id, nil)
		return
	}

	completions := s.getCompletions()
	s.sendResponse(id, map[string]interface{}{
		"isIncomplete": false,
		"items":        completions,
	})
}

func (s *LSPServer) handleHover(id interface{}, params json.RawMessage) {
	var request struct {
		TextDocument struct {
			URI string `json:"uri"`
		} `json:"textDocument"`
		Position struct {
			Line      int `json:"line"`
			Character int `json:"character"`
		} `json:"position"`
	}

	if err := json.Unmarshal(params, &request); err != nil {
		s.sendResponse(id, nil)
		return
	}

	s.mu.Lock()
	text := s.documents[request.TextDocument.URI]
	s.mu.Unlock()

	hover := s.getHover(text, request.Position.Line, request.Position.Character)
	s.sendResponse(id, hover)
}

func (s *LSPServer) sendResponse(id interface{}, result interface{}) {
	response := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"result":  result,
	}

	content, _ := json.Marshal(response)
	s.writeMessage(content)
}

func (s *LSPServer) sendNotification(method string, params interface{}) {
	notification := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
	}

	content, _ := json.Marshal(notification)
	s.writeMessage(content)
}

func (s *LSPServer) publishDiagnostics(uri, text string) {
	diagnostics := s.getDiagnostics(text)

	params := map[string]interface{}{
		"uri":         uri,
		"diagnostics": diagnostics,
		"version":     0,
	}

	s.sendNotification("textDocument/publishDiagnostics", params)
}

func (s *LSPServer) getDiagnostics(text string) []interface{} {
	var diagnostics []interface{}

	// Try to compile and get errors
	c := compiler.New("ink")
	_, err := c.Compile("dummy.gox", text)
	if err != nil {
		// Parse error message for line/column
		diagnostic := map[string]interface{}{
			"range": map[string]interface{}{
				"start": map[string]int{"line": 0, "character": 0},
				"end":   map[string]int{"line": 0, "character": 1},
			},
			"severity": 1, // Error
			"source":   "gox",
			"message":  err.Error(),
		}
		diagnostics = append(diagnostics, diagnostic)
	}

	return diagnostics
}

func (s *LSPServer) getCompletions() []interface{} {
	// JSX component completions
	components := []string{"Box", "Text", "Spacer", "Newline", "Static", "Transform", "Fragment"}
	hooks := []string{"useState", "useEffect", "useLayoutEffect", "useRef", "useMemo", "useCallback", "useInput", "useApp", "useFocus", "useFocusManager", "useCursor", "useAnimation"}

	var items []interface{}

	for _, comp := range components {
		items = append(items, map[string]interface{}{
			"label":            comp,
			"kind":             7, // Class
			"detail":           fmt.Sprintf("Component: %s", comp),
			"insertText":       fmt.Sprintf("<%s>$1</%s>", comp, comp),
			"insertTextFormat": 2, // Snippet
		})
	}

	for _, hook := range hooks {
		items = append(items, map[string]interface{}{
			"label":  hook,
			"kind":   3, // Function
			"detail": fmt.Sprintf("Hook: %s", hook),
		})
	}

	return items
}

func (s *LSPServer) getHover(text string, line, character int) interface{} {
	// Find word at position
	lines := strings.Split(text, "\n")
	if line >= len(lines) {
		return nil
	}

	lineText := lines[line]
	if character > len(lineText) {
		return nil
	}

	// Extract word
	start := character
	for start > 0 && isWordChar(lineText[start-1]) {
		start--
	}
	end := character
	for end < len(lineText) && isWordChar(lineText[end]) {
		end++
	}

	if start == end {
		return nil
	}

	word := lineText[start:end]
	docs := getDocumentation(word)

	if docs == "" {
		return nil
	}

	return map[string]interface{}{
		"contents": map[string]interface{}{
			"kind":  "markdown",
			"value": docs,
		},
	}
}

func isWordChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_'
}

func getDocumentation(word string) string {
	// Component documentation
	componentDocs := map[string]string{
		"Box": "## Box Component\n\nA flexbox container component.\n\n**Props:**\n- `width`: Width (number or string)\n- `height`: Height (number or string)\n- `flexDirection`: \"row\" | \"column\" | \"row-reverse\" | \"column-reverse\"\n- `justifyContent`: \"flex-start\" | \"center\" | \"flex-end\" | \"space-between\" | \"space-around\"\n- `alignItems`: \"flex-start\" | \"center\" | \"flex-end\" | \"stretch\"\n- `padding`, `paddingX`, `paddingY`: Padding\n- `margin`, `marginX`, `marginY`: Margin\n- `borderStyle`: \"single\" | \"double\" | \"round\" | \"bold\"\n- `borderColor`: Border color",

		"Text": "## Text Component\n\nText content with styling support.\n\n**Props:**\n- `color`: Text color\n- `backgroundColor`: Background color\n- `bold`: Bold text\n- `italic`: Italic text\n- `underline`: Underlined text\n- `dimColor`: Dimmed text\n- `inverse`: Invert colors\n- `wrap`: \"wrap\" | \"truncate\" | \"truncate-start\" | \"truncate-middle\" | \"truncate-end\"",

		"Spacer": "## Spacer Component\n\nFlexible whitespace that expands to fill available space.\n\n**Props:**\n- `width`: Minimum width\n- `height`: Minimum height",

		"Newline": "## Newline Component\n\nLine break component.\n\n**Props:**\n- `count`: Number of newlines (default: 1)",
	}

	// Hook documentation
	hookDocs := map[string]string{
		"useState": "## useState Hook\n\nState management hook.\n\n```go\ncount, setCount := hooks.UseState(0)\nsetCount(count + 1)\n```",

		"useEffect": "## useEffect Hook\n\nSide effect hook.\n\n```go\nhooks.UseEffect(func() func() {\n    // Setup\n    return func() {\n        // Cleanup\n    }\n}, []interface{}{deps})\n```",

		"useRef": "## useRef Hook\n\nMutable reference hook.\n\n```go\nref := hooks.UseRef(initialValue)\nref.Current = newValue\n```",

		"useMemo": "## useMemo Hook\n\nMemoization hook.\n\n```go\nvalue := hooks.UseMemo(func() interface{} {\n    return expensiveComputation()\n}, []interface{}{deps})\n```",

		"useInput": "## useInput Hook\n\nKeyboard input hook.\n\n```go\nhooks.UseInput(func(key string, mod hooks.ModMask) {\n    if key == \"q\" {\n        app.Quit()\n    }\n})\n```",

		"useApp": "## useApp Hook\n\nApplication context hook.\n\n```go\napp := hooks.UseApp()\napp.Quit()\napp.Rerender()\n```",

		"useFocus": "## useFocus Hook\n\nFocus management hook.\n\n```go\nfocus, setFocus := hooks.UseFocus(true)\nif focus {\n    // Render focused state\n}\n```",

		"useCursor": "## useCursor Hook\n\nCursor visibility hook.\n\n```go\nshow, hide := hooks.UseCursor()\nshow()\nhide()\n```",

		"useAnimation": "## useAnimation Hook\n\nAnimation hook.\n\n```go\nframe, start, stop := hooks.UseAnimation(100 * time.Millisecond)\nstart()\n// Use frame value\nstop()\n```",
	}

	if doc, ok := componentDocs[word]; ok {
		return doc
	}
	if doc, ok := hookDocs[word]; ok {
		return doc
	}
	return ""
}
