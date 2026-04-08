package com.ayanmw.gox

import com.intellij.openapi.editor.colors.TextAttributesKey
import com.intellij.openapi.fileTypes.SyntaxHighlighter
import com.intellij.openapi.fileTypes.SyntaxHighlighterFactory
import com.intellij.openapi.project.Project
import com.intellij.openapi.vfs.VirtualFile
import com.intellij.openapi.editor.DefaultLanguageHighlighterColors as Defaults

/**
 * Syntax highlighter factory for GoX
 */
class GoxSyntaxHighlighterFactory : SyntaxHighlighterFactory() {
    override fun getSyntaxHighlighter(project: Project?, virtualFile: VirtualFile?): SyntaxHighlighter {
        return GoxSyntaxHighlighter()
    }
}

/**
 * Syntax highlighter for GoX language
 */
class GoxSyntaxHighlighter : SyntaxHighlighter {
    companion object {
        // Go keywords
        val KEYWORD = TextAttributesKey.createTextAttributesKey("GOX_KEYWORD", Defaults.KEYWORD)
        val KEYWORD_TYPE = TextAttributesKey.createTextAttributesKey("GOX_KEYWORD_TYPE", Defaults.KEYWORD)
        val BUILTIN_FUNC = TextAttributesKey.createTextAttributesKey("GOX_BUILTIN_FUNC", Defaults.PREDEFINED_SYMBOL)

        // Literals
        val STRING = TextAttributesKey.createTextAttributesKey("GOX_STRING", Defaults.STRING)
        val NUMBER = TextAttributesKey.createTextAttributesKey("GOX_NUMBER", Defaults.NUMBER)
        val CONSTANT = TextAttributesKey.createTextAttributesKey("GOX_CONSTANT", Defaults.CONSTANT)

        // Comments
        val LINE_COMMENT = TextAttributesKey.createTextAttributesKey("GOX_LINE_COMMENT", Defaults.LINE_COMMENT)
        val BLOCK_COMMENT = TextAttributesKey.createTextAttributesKey("GOX_BLOCK_COMMENT", Defaults.BLOCK_COMMENT)

        // Operators
        val OPERATOR = TextAttributesKey.createTextAttributesKey("GOX_OPERATOR", Defaults.OPERATION_SIGN)
        val BRACKET = TextAttributesKey.createTextAttributesKey("GOX_BRACKET", Defaults.BRACKETS)
        val PAREN = TextAttributesKey.createTextAttributesKey("GOX_PAREN", Defaults.PARENTHESES)
        val BRACE = TextAttributesKey.createTextAttributesKey("GOX_BRACE", Defaults.BRACES)

        // JSX
        val JSX_TAG = TextAttributesKey.createTextAttributesKey("GOX_JSX_TAG", Defaults.MARKUP_TAG)
        val JSX_TAG_NAME = TextAttributesKey.createTextAttributesKey("GOX_JSX_TAG_NAME", Defaults.MARKUP_TAG_NAME)
        val JSX_ATTR_NAME = TextAttributesKey.createTextAttributesKey("GOX_JSX_ATTR_NAME", Defaults.MARKUP_ATTRIBUTE)
        val JSX_ATTR_VALUE = TextAttributesKey.createTextAttributesKey("GOX_JSX_ATTR_VALUE", Defaults.MARKUP_ENTITY)
        val JSX_EXPR = TextAttributesKey.createTextAttributesKey("GOX_JSX_EXPR", Defaults.TEMPLATE_LANGUAGE_COLOR)

        // Identifiers
        val IDENTIFIER = TextAttributesKey.createTextAttributesKey("GOX_IDENTIFIER", Defaults.IDENTIFIER)
        val FUNCTION_NAME = TextAttributesKey.createTextAttributesKey("GOX_FUNCTION_NAME", Defaults.FUNCTION_DECLARATION)

        // Keywords
        val KEYWORD_TOKENS = arrayOf(
            GoxTokenTypes.PACKAGE, GoxTokenTypes.IMPORT, GoxTokenTypes.FUNC,
            GoxTokenTypes.RETURN, GoxTokenTypes.VAR, GoxTokenTypes.CONST,
            GoxTokenTypes.TYPE, GoxTokenTypes.STRUCT, GoxTokenTypes.INTERFACE,
            GoxTokenTypes.MAP, GoxTokenTypes.CHAN,
            GoxTokenTypes.IF, GoxTokenTypes.ELSE,
            GoxTokenTypes.FOR, GoxTokenTypes.RANGE,
            GoxTokenTypes.SWITCH, GoxTokenTypes.CASE, GoxTokenTypes.DEFAULT,
            GoxTokenTypes.SELECT, GoxTokenTypes.GO, GoxTokenTypes.DEFER
        )

        val OPERATOR_TOKENS = arrayOf(
            GoxTokenTypes.PLUS, GoxTokenTypes.MINUS, GoxTokenTypes.STAR,
            GoxTokenTypes.SLASH, GoxTokenTypes.PERCENT,
            GoxTokenTypes.AND, GoxTokenTypes.OR, GoxTokenTypes.XOR,
            GoxTokenTypes.LAND, GoxTokenTypes.LOR, GoxTokenTypes.NOT,
            GoxTokenTypes.ASSIGN, GoxTokenTypes.DEFINE,
            GoxTokenTypes.EQ, GoxTokenTypes.NE,
            GoxTokenTypes.LT, GoxTokenTypes.LE, GoxTokenTypes.GT, GoxTokenTypes.GE
        )
    }

    override fun getHighlightingLexer() = GoxLexerAdapter()
    override fun getTokenHighlights(tokenType: com.intellij.psi.tree.IElementType): Array<TextAttributesKey> {
        return when (tokenType) {
            in KEYWORD_TOKENS -> arrayOf(KEYWORD)
            GoxTokenTypes.STRING, GoxTokenTypes.RAW_STRING -> arrayOf(STRING)
            GoxTokenTypes.CHAR -> arrayOf(STRING)
            GoxTokenTypes.NUMBER -> arrayOf(NUMBER)
            GoxTokenTypes.LINE_COMMENT -> arrayOf(LINE_COMMENT)
            GoxTokenTypes.BLOCK_COMMENT -> arrayOf(BLOCK_COMMENT)
            in OPERATOR_TOKENS -> arrayOf(OPERATOR)
            GoxTokenTypes.LPAREN, GoxTokenTypes.RPAREN -> arrayOf(PAREN)
            GoxTokenTypes.LBRACE, GoxTokenTypes.RBRACE -> arrayOf(BRACE)
            GoxTokenTypes.LBRACK, GoxTokenTypes.RBRACK -> arrayOf(BRACKET)
            GoxTokenTypes.JSX_TAG_NAME -> arrayOf(JSX_TAG_NAME)
            GoxTokenTypes.JSX_ATTR_NAME -> arrayOf(JSX_ATTR_NAME)
            GoxTokenTypes.JSX_ATTR_VALUE -> arrayOf(JSX_ATTR_VALUE)
            GoxTokenTypes.JSX_TAG_OPEN, GoxTokenTypes.JSX_TAG_CLOSE -> arrayOf(JSX_TAG)
            GoxTokenTypes.JSX_EXPR_START, GoxTokenTypes.JSX_EXPR_END -> arrayOf(JSX_EXPR)
            else -> emptyArray()
        }
    }
}