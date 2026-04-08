package com.ayanmw.gox

import com.intellij.lang.ASTNode
import com.intellij.lang.PsiBuilder
import com.intellij.lang.PsiParser
import com.intellij.psi.tree.IElementType

/**
 * Parser for GoX language
 */
class GoxParser : PsiParser {
    override fun parse(root: IElementType, builder: PsiBuilder): ASTNode {
        val rootMarker = builder.mark()

        while (!builder.eof()) {
            val tokenType = builder.tokenType

            when (tokenType) {
                GoxTokenTypes.JSX_TAG_NAME -> parseJSXElement(builder)
                GoxTokenTypes.PACKAGE -> parsePackage(builder)
                GoxTokenTypes.IMPORT -> parseImport(builder)
                GoxTokenTypes.FUNC -> parseFunc(builder)
                else -> builder.advanceLexer()
            }
        }

        rootMarker.done(root)
        return builder.treeBuilt
    }

    private fun parseJSXElement(builder: PsiBuilder) {
        val marker = builder.mark()
        builder.advanceLexer() // tag name

        // Parse attributes
        while (builder.tokenType == GoxTokenTypes.IDENTIFIER ||
               builder.tokenType == GoxTokenTypes.JSX_ATTR_NAME) {
            builder.advanceLexer() // attr name
            if (builder.tokenType == GoxTokenTypes.ASSIGN) {
                builder.advanceLexer() // =
                if (builder.tokenType == GoxTokenTypes.STRING) {
                    builder.advanceLexer() // string value
                } else if (builder.tokenType == GoxTokenTypes.LBRACE) {
                    parseExpression(builder)
                }
            }
        }

        // Check for self-closing or children
        if (builder.tokenType == GoxTokenTypes.SLASH) {
            builder.advanceLexer() // /
            builder.advanceLexer() // >
            marker.done(GoxTokenTypes.JSX_TAG_SELF_CLOSE)
        } else if (builder.tokenType == GoxTokenTypes.GT) {
            builder.advanceLexer() // >
            // Parse children
            while (builder.tokenType != GoxTokenTypes.LT && !builder.eof()) {
                when (builder.tokenType) {
                    GoxTokenTypes.JSX_TAG_NAME -> parseJSXElement(builder)
                    GoxTokenTypes.JSX_EXPR_START -> parseExpression(builder)
                    else -> builder.advanceLexer()
                }
            }
            // Closing tag
            if (builder.tokenType == GoxTokenTypes.LT) {
                builder.advanceLexer() // <
                if (builder.tokenType == GoxTokenTypes.SLASH) {
                    builder.advanceLexer() // /
                    builder.advanceLexer() // tag name
                    builder.advanceLexer() // >
                }
            }
            marker.done(GoxTokenTypes.JSX_TAG_CLOSE)
        } else {
            marker.drop()
        }
    }

    private fun parseExpression(builder: PsiBuilder) {
        val marker = builder.mark()
        builder.advanceLexer() // {

        var depth = 1
        while (depth > 0 && !builder.eof()) {
            when (builder.tokenType) {
                GoxTokenTypes.LBRACE -> depth++
                GoxTokenTypes.RBRACE -> depth--
            }
            builder.advanceLexer()
        }

        marker.done(GoxTokenTypes.JSX_EXPR_END)
    }

    private fun parsePackage(builder: PsiBuilder) {
        val marker = builder.mark()
        builder.advanceLexer() // package
        builder.advanceLexer() // name
        marker.done(GoxTokenTypes.PACKAGE)
    }

    private fun parseImport(builder: PsiBuilder) {
        val marker = builder.mark()
        builder.advanceLexer() // import
        if (builder.tokenType == GoxTokenTypes.LPAREN) {
            builder.advanceLexer() // (
            while (builder.tokenType != GoxTokenTypes.RPAREN && !builder.eof()) {
                builder.advanceLexer()
            }
            builder.advanceLexer() // )
        } else {
            builder.advanceLexer() // import spec
        }
        marker.done(GoxTokenTypes.IMPORT)
    }

    private fun parseFunc(builder: PsiBuilder) {
        val marker = builder.mark()
        builder.advanceLexer() // func
        builder.advanceLexer() // name
        if (builder.tokenType == GoxTokenTypes.LPAREN) {
            parseParameters(builder)
        }
        // Skip to function body
        while (builder.tokenType != GoxTokenTypes.LBRACE && !builder.eof()) {
            builder.advanceLexer()
        }
        if (builder.tokenType == GoxTokenTypes.LBRACE) {
            parseBlock(builder)
        }
        marker.done(GoxTokenTypes.FUNC)
    }

    private fun parseParameters(builder: PsiBuilder) {
        builder.advanceLexer() // (
        while (builder.tokenType != GoxTokenTypes.RPAREN && !builder.eof()) {
            builder.advanceLexer()
        }
        builder.advanceLexer() // )
    }

    private fun parseBlock(builder: PsiBuilder) {
        builder.advanceLexer() // {
        var depth = 1
        while (depth > 0 && !builder.eof()) {
            when (builder.tokenType) {
                GoxTokenTypes.LBRACE -> depth++
                GoxTokenTypes.RBRACE -> depth--
            }
            builder.advanceLexer()
        }
    }
}