package com.ayanmw.gox

import com.intellij.lexer.FlexAdapter
import com.intellij.lexer.Lexer
import com.intellij.openapi.project.Project
import java.io.Reader

/**
 * Lexer adapter for GoX
 */
class GoxLexerAdapter : FlexAdapter(GoxLexer(null as? Reader?))

/**
 * Flex-based lexer for GoX language
 * Handles both Go code and JSX syntax
 */
class GoxLexer(reader: Reader?) : com.intellij.lexer.FlexLexer {
    private var buffer: CharSequence = ""
    private var startOffset: Int = 0
    private var endOffset: Int = 0
    private var position: Int = 0
    private var tokenStart: Int = 0
    private var tokenEnd: Int = 0
    private var currentToken: IElementType? = null

    // State tracking
    private var inJSX: Boolean = false
    private var jsxDepth: Int = 0

    override fun yybegin(newState: Int) {}
    override fun yystate(): Int = if (inJSX) 1 else 0
    override fun getBufferEnd(): Int = endOffset
    override fun getBufferSequence(): CharSequence = buffer

    override fun advance(): IElementType? {
        if (position >= endOffset) {
            currentToken = null
            return null
        }

        tokenStart = position
        val ch = buffer[position]

        // Check for JSX start
        if (ch == '<' && position + 1 < endOffset) {
            val nextCh = buffer[position + 1]
            if (nextCh.isUpperCase() || nextCh == '>') {
                return parseJSXTag()
            }
        }

        // Check for JSX expression
        if (inJSX && ch == '{') {
            return parseJSXExpression()
        }

        // Parse Go code
        return parseGoCode()
    }

    private fun parseJSXTag(): IElementType {
        val ch = buffer[position + 1]

        if (ch == '>') {
            // Fragment <>
            position += 2
            tokenEnd = position
            jsxDepth++
            inJSX = true
            return GoxTokenTypes.JSX_TAG_OPEN
        }

        // Find tag name
        position++ // skip <
        val tagNameStart = position
        while (position < endOffset && (buffer[position].isLetterOrDigit() || buffer[position] == '_')) {
            position++
        }

        tokenEnd = position
        inJSX = true
        jsxDepth++
        return GoxTokenTypes.JSX_TAG_NAME
    }

    private fun parseJSXExpression(): IElementType {
        position++ // skip {
        tokenEnd = position
        return GoxTokenTypes.JSX_EXPR_START
    }

    private fun parseGoCode(): IElementType? {
        val ch = buffer[position]

        when {
            ch.isWhitespace() -> {
                while (position < endOffset && buffer[position].isWhitespace()) {
                    position++
                }
                tokenEnd = position
                return GoxTokenTypes.WHITE_SPACE
            }
            ch == '/' && position + 1 < endOffset -> {
                if (buffer[position + 1] == '/') {
                    // Line comment
                    while (position < endOffset && buffer[position] != '\n') {
                        position++
                    }
                    tokenEnd = position
                    return GoxTokenTypes.LINE_COMMENT
                } else if (buffer[position + 1] == '*') {
                    // Block comment
                    position += 2
                    while (position + 1 < endOffset && !(buffer[position] == '*' && buffer[position + 1] == '/')) {
                        position++
                    }
                    position += 2
                    tokenEnd = position
                    return GoxTokenTypes.BLOCK_COMMENT
                }
            }
            ch == '"' -> {
                position++
                while (position < endOffset && buffer[position] != '"') {
                    if (buffer[position] == '\\') position++
                    position++
                }
                position++
                tokenEnd = position
                return GoxTokenTypes.STRING
            }
            ch == '`' -> {
                position++
                while (position < endOffset && buffer[position] != '`') {
                    position++
                }
                position++
                tokenEnd = position
                return GoxTokenTypes.RAW_STRING
            }
            ch == '\'' -> {
                position++
                while (position < endOffset && buffer[position] != '\'') {
                    if (buffer[position] == '\\') position++
                    position++
                }
                position++
                tokenEnd = position
                return GoxTokenTypes.CHAR
            }
            ch.isDigit() -> {
                while (position < endOffset && (buffer[position].isDigit() || buffer[position] in "xXoObBeE.+-")) {
                    position++
                }
                tokenEnd = position
                return GoxTokenTypes.NUMBER
            }
            ch.isLetter() || ch == '_' -> {
                val start = position
                while (position < endOffset && (buffer[position].isLetterOrDigit() || buffer[position] == '_')) {
                    position++
                }
                tokenEnd = position
                val text = buffer.subSequence(start, position).toString()
                return GoxTokenTypes.KEYWORDS[text] ?: GoxTokenTypes.IDENTIFIER
            }
            else -> {
                position++
                tokenEnd = position
                return when (ch) {
                    '+' -> GoxTokenTypes.PLUS
                    '-' -> GoxTokenTypes.MINUS
                    '*' -> GoxTokenTypes.STAR
                    '/' -> GoxTokenTypes.SLASH
                    '%' -> GoxTokenTypes.PERCENT
                    '&' -> GoxTokenTypes.AND
                    '|' -> GoxTokenTypes.OR
                    '^' -> GoxTokenTypes.XOR
                    '!' -> GoxTokenTypes.NOT
                    '=' -> if (position < endOffset && buffer[position] == '=') {
                        position++; tokenEnd = position; GoxTokenTypes.EQ
                    } else GoxTokenTypes.ASSIGN
                    '<' -> if (position < endOffset && buffer[position] == '=') {
                        position++; tokenEnd = position; GoxTokenTypes.LE
                    } else GoxTokenTypes.LT
                    '>' -> if (position < endOffset && buffer[position] == '=') {
                        position++; tokenEnd = position; GoxTokenTypes.GE
                    } else GoxTokenTypes.GT
                    '(' -> GoxTokenTypes.LPAREN
                    ')' -> GoxTokenTypes.RPAREN
                    '{' -> GoxTokenTypes.LBRACE
                    '}' -> GoxTokenTypes.RBRACE
                    '[' -> GoxTokenTypes.LBRACK
                    ']' -> GoxTokenTypes.RBRACK
                    ',' -> GoxTokenTypes.COMMA
                    ';' -> GoxTokenTypes.SEMICOLON
                    '.' -> GoxTokenTypes.DOT
                    ':' -> if (position < endOffset && buffer[position] == '=') {
                        position++; tokenEnd = position; GoxTokenTypes.DEFINE
                    } else null
                    else -> null
                }
            }
        }

        return null
    }

    override fun getTokenStart(): Int = tokenStart
    override fun getTokenEnd(): Int = tokenEnd
    override fun getTokenType(): IElementType? = currentToken

    override fun reset(buffer: CharSequence, startOffset: Int, endOffset: Int, initialState: Int) {
        this.buffer = buffer
        this.startOffset = startOffset
        this.endOffset = endOffset
        this.position = startOffset
        this.tokenStart = startOffset
        this.tokenEnd = startOffset
        this.currentToken = null
        this.inJSX = false
        this.jsxDepth = 0
    }
}

typealias IElementType = com.intellij.psi.tree.IElementType