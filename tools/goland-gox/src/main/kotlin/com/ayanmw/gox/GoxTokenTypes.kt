package com.ayanmw.gox

import com.intellij.psi.tree.IElementType

/**
 * Token types for GoX language
 */
class GoxTokenType(debugName: String) : IElementType(debugName, GoxLanguage) {
    override fun toString(): String = "GoxTokenType." + super.toString()
}

/**
 * Token types constants
 */
object GoxTokenTypes {
    // Whitespace and Comments
    val WHITE_SPACE = GoxTokenType("WHITE_SPACE")
    val LINE_COMMENT = GoxTokenType("LINE_COMMENT")
    val BLOCK_COMMENT = GoxTokenType("BLOCK_COMMENT")

    // Keywords
    val PACKAGE = GoxTokenType("package")
    val IMPORT = GoxTokenType("import")
    val FUNC = GoxTokenType("func")
    val RETURN = GoxTokenType("return")
    val VAR = GoxTokenType("var")
    val CONST = GoxTokenType("const")
    val TYPE = GoxTokenType("type")
    val STRUCT = GoxTokenType("struct")
    val INTERFACE = GoxTokenType("interface")
    val MAP = GoxTokenType("map")
    val CHAN = GoxTokenType("chan")
    val IF = GoxTokenType("if")
    val ELSE = GoxTokenType("else")
    val FOR = GoxTokenType("for")
    val RANGE = GoxTokenType("range")
    val SWITCH = GoxTokenType("switch")
    val CASE = GoxTokenType("case")
    val DEFAULT = GoxTokenType("default")
    val SELECT = GoxTokenType("select")
    val GO = GoxTokenType("go")
    val DEFER = GoxTokenType("defer")

    // Literals
    val IDENTIFIER = GoxTokenType("IDENTIFIER")
    val STRING = GoxTokenType("STRING")
    val RAW_STRING = GoxTokenType("RAW_STRING")
    val CHAR = GoxTokenType("CHAR")
    val NUMBER = GoxTokenType("NUMBER")

    // Operators
    val PLUS = GoxTokenType("+")
    val MINUS = GoxTokenType("-")
    val STAR = GoxTokenType("*")
    val SLASH = GoxTokenType("/")
    val PERCENT = GoxTokenType("%")
    val AND = GoxTokenType("&")
    val OR = GoxTokenType("|")
    val XOR = GoxTokenType("^")
    val LAND = GoxTokenType("&&")
    val LOR = GoxTokenType("||")
    val NOT = GoxTokenType("!")
    val ASSIGN = GoxTokenType("=")
    val DEFINE = GoxTokenType(":=")
    val EQ = GoxTokenType("==")
    val NE = GoxTokenType("!=")
    val LT = GoxTokenType("<")
    val LE = GoxTokenType("<=")
    val GT = GoxTokenType(">")
    val GE = GoxTokenType(">=")

    // Delimiters
    val LPAREN = GoxTokenType("(")
    val RPAREN = GoxTokenType(")")
    val LBRACE = GoxTokenType("{")
    val RBRACE = GoxTokenType("}")
    val LBRACK = GoxTokenType("[")
    val RBRACK = GoxTokenType("]")
    val COMMA = GoxTokenType(",")
    val SEMICOLON = GoxTokenType(";")
    val DOT = GoxTokenType(".")

    // JSX Tokens
    val JSX_TAG_OPEN = GoxTokenType("JSX_TAG_OPEN")
    val JSX_TAG_CLOSE = GoxTokenType("JSX_TAG_CLOSE")
    val JSX_TAG_SELF_CLOSE = GoxTokenType("JSX_TAG_SELF_CLOSE")
    val JSX_TAG_NAME = GoxTokenType("JSX_TAG_NAME")
    val JSX_ATTR_NAME = GoxTokenType("JSX_ATTR_NAME")
    val JSX_ATTR_VALUE = GoxTokenType("JSX_ATTR_VALUE")
    val JSX_TEXT = GoxTokenType("JSX_TEXT")
    val JSX_EXPR_START = GoxTokenType("JSX_EXPR_START")
    val JSX_EXPR_END = GoxTokenType("JSX_EXPR_END")

    // Keywords map
    val KEYWORDS = mapOf(
        "package" to PACKAGE,
        "import" to IMPORT,
        "func" to FUNC,
        "return" to RETURN,
        "var" to VAR,
        "const" to CONST,
        "type" to TYPE,
        "struct" to STRUCT,
        "interface" to INTERFACE,
        "map" to MAP,
        "chan" to CHAN,
        "if" to IF,
        "else" to ELSE,
        "for" to FOR,
        "range" to RANGE,
        "switch" to SWITCH,
        "case" to CASE,
        "default" to DEFAULT,
        "select" to SELECT,
        "go" to GO,
        "defer" to DEFER
    )

    object Factory {
        fun createElement(node: com.intellij.lang.ASTNode): com.intellij.psi.PsiElement {
            return com.intellij.psi.impl.source.tree.LeafPsiElementImpl(node.elementType, node.text)
        }
    }
}