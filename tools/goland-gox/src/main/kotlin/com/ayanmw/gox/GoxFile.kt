package com.ayanmw.gox

import com.intellij.psi.FileViewProvider
import com.intellij.psi.PsiFileBase

/**
 * GoX PSI File
 */
class GoxFile(viewProvider: FileViewProvider) : PsiFileBase(viewProvider, GoxLanguage) {
    override fun getFileType() = GoxFileType.INSTANCE
    override fun toString() = "GoX File"
}