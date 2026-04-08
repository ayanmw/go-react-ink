package com.ayanmw.gox

import com.intellij.lang.Language

/**
 * GoX Language definition for JSX in Go
 */
object GoxLanguage : Language("GoX") {
    override fun getDisplayName(): String = "GoX"
    override fun getMimeType(): String = "text/x-gox"
    override fun getAssociatedFileType(): GoxFileType = GoxFileType.INSTANCE
}