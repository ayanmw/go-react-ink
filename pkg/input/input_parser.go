package input

import (
	"strings"
)

// ParsedKey 解析后的按键信息（用于 ParseKey）
type ParsedKey struct {
	UpArrow    bool
	DownArrow  bool
	LeftArrow  bool
	RightArrow bool
	PageDown   bool
	PageUp     bool
	Home       bool
	End        bool
	Return     bool
	Escape     bool
	Ctrl       bool
	Shift      bool
	Tab        bool
	Backspace  bool
	Delete     bool
	Meta       bool
}

type Keypress struct {
	Name     string
	Sequence string
	Ctrl     bool
	Shift    bool
	Meta     bool
}

var arrowKeys = map[string]bool{
	"up": true, "down": true, "left": true, "right": true,
}

var specialKeys = map[string]bool{
	"pagedown": true, "pageup": true, "home": true, "end": true,
	"return": true, "escape": true, "tab": true,
	"backspace": true, "delete": true,
}

// ParseKey 解析 ANSI 键盘序列，返回原有的 Key 类型
func ParseKey(data string) Key {
	parsedKey, sequence := parseKeypressToKey(data)

	// 转换为原有的 Key 类型
	key := Key{
		Sequence: sequence,
		Ctrl:     parsedKey.Ctrl,
		Shift:    parsedKey.Shift,
		Alt:      parsedKey.Meta, // 映射 Meta 到 Alt
		Meta:     parsedKey.Meta,
	}
	return key
}

func parseKeypressToKey(data string) (ParsedKey, string) {
	kp := parseKeypress(data)
	if kp == nil {
		return ParsedKey{}, ""
	}

	key := ParsedKey{}
	name := strings.ToLower(kp.Name)

	if arrowKeys[name] {
		switch name {
		case "up":
			key.UpArrow = true
		case "down":
			key.DownArrow = true
		case "left":
			key.LeftArrow = true
		case "right":
			key.RightArrow = true
		}
	}

	if specialKeys[name] {
		switch name {
		case "pagedown":
			key.PageDown = true
		case "pageup":
			key.PageUp = true
		case "home":
			key.Home = true
		case "end":
			key.End = true
		case "return":
			key.Return = true
		case "escape":
			key.Escape = true
		case "tab":
			key.Tab = true
		case "backspace":
			key.Backspace = true
		case "delete":
			key.Delete = true
		}
	}

	key.Ctrl = kp.Ctrl
	key.Shift = kp.Shift
	key.Meta = kp.Meta
	key.Return = name == "return"

	return key, kp.Sequence
}

// 简化版的 parseKeypress，核心功能
func parseKeypress(data string) *Keypress {
	if len(data) == 0 {
		return nil
	}

	kp := &Keypress{}

	// 处理 Ctrl 组合键 (Ctrl+C = \x03)
	if len(data) == 1 && data[0] >= 1 && data[0] <= 26 {
		kp.Ctrl = true
		kp.Name = string(rune(data[0] + 'a' - 1))
		kp.Sequence = data
		return kp
	}

	// 处理 ANSI 转义序列
	if strings.HasPrefix(data, "\x1b") {
		return parseEscapeSequence(data)
	}

	// 普通字符
	kp.Name = data
	kp.Sequence = data
	return kp
}

// 解析 ANSI 转义序列（简化版，仅覆盖常用）
func parseEscapeSequence(data string) *Keypress {
	kp := &Keypress{}

	if data == "\x1b" {
		kp.Name = "escape"
		kp.Sequence = data
		return kp
	}

	if strings.HasPrefix(data, "\x1b[") {
		rest := data[2:]
		if strings.HasSuffix(rest, "~") {
			// 处理特殊序列
			num := strings.TrimSuffix(rest, "~")
			switch num {
			case "3":
				kp.Name = "delete"
			case "5":
				kp.Name = "pageup"
			case "6":
				kp.Name = "pagedown"
			case "H":
				kp.Name = "home"
			case "F":
				kp.Name = "end"
			}
		} else {
			// 处理箭头键和其他
			if strings.HasSuffix(rest, "A") {
				kp.Name = "up"
			} else if strings.HasSuffix(rest, "B") {
				kp.Name = "down"
			} else if strings.HasSuffix(rest, "C") {
				kp.Name = "right"
			} else if strings.HasSuffix(rest, "D") {
				kp.Name = "left"
			}
		}
	} else if strings.HasPrefix(data, "\x1b\t") {
		kp.Name = "tab"
		kp.Shift = true
	} else if strings.HasPrefix(data, "\x1b[") {
		// 处理 [1;5A 等格式（Ctrl+箭头）
		rest := data[2:]
		if strings.HasSuffix(rest, "A") {
			kp.Name = "up"
			if strings.Contains(rest, "5") || strings.Contains(rest, "6") {
				kp.Ctrl = true
			}
		} else if strings.HasSuffix(rest, "B") {
			kp.Name = "down"
			if strings.Contains(rest, "5") || strings.Contains(rest, "6") {
				kp.Ctrl = true
			}
		} else if strings.HasSuffix(rest, "C") {
			kp.Name = "right"
			if strings.Contains(rest, "5") || strings.Contains(rest, "6") {
				kp.Ctrl = true
			}
		} else if strings.HasSuffix(rest, "D") {
			kp.Name = "left"
			if strings.Contains(rest, "5") || strings.Contains(rest, "6") {
				kp.Ctrl = true
			}
		}
	}

	kp.Sequence = data
	return kp
}
