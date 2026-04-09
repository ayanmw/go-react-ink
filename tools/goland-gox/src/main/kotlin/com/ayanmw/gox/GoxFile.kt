package com.ayanmw.gox

import com.intellij.extapi.psi.PsiFileBase
import com.intellij.psi.FileViewProvider

/**
 * GoX PSI File
 */
class GoxFile(viewProvider: FileViewProvider) : PsiFileBase(viewProvider, GoxLanguage) {
    override fun getFileType() = GoxFileType.INSTANCE
    override fun toString() = "GoX File"
}