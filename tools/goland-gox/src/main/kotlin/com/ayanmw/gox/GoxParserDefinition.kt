package com.ayanmw.gox

import com.intellij.lang.ASTNode
import com.intellij.lang.ParserDefinition
import com.intellij.lang.PsiParser
import com.intellij.lexer.Lexer
import com.intellij.openapi.project.Project
import com.intellij.psi.FileViewProvider
import com.intellij.psi.PsiElement
import com.intellij.psi.PsiFile
import com.intellij.psi.tree.IFileElementType
import com.intellij.psi.tree.TokenSet

/**
 * Parser definition for GoX language
 */
class GoxParserDefinition : ParserDefinition {
    companion object {
        val FILE = IFileElementType(GoxLanguage)
        val WHITE_SPACES = TokenSet.create(GoxTokenTypes.WHITE_SPACE)
        val COMMENTS = TokenSet.create(GoxTokenTypes.LINE_COMMENT, GoxTokenTypes.BLOCK_COMMENT)
        val STRING_LITERALS = TokenSet.create(GoxTokenTypes.STRING, GoxTokenTypes.RAW_STRING, GoxTokenTypes.CHAR)
    }

    override fun createLexer(project: Project?): Lexer = GoxLexerAdapter()
    override fun createParser(project: Project?): PsiParser = GoxParser()
    override fun getFileNodeType(): IFileElementType = FILE
    override fun getCommentTokens(): TokenSet = COMMENTS
    override fun getStringLiteralElements(): TokenSet = STRING_LITERALS

    override fun createElement(node: ASTNode): PsiElement = GoxTokenTypes.Factory.createElement(node)
    override fun createFile(viewProvider: FileViewProvider): PsiFile = GoxFile(viewProvider)
}