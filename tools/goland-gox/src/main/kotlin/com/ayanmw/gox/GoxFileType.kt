package com.ayanmw.gox

import com.intellij.openapi.fileTypes.LanguageFileType
import javax.swing.Icon

/**
 * GoX File Type definition
 */
object GoxFileType : LanguageFileType(GoxLanguage) {
    val INSTANCE = this

    override fun getName(): String = "GoX"
    override fun getDescription(): String = "GoX - JSX for Go"
    override fun getDefaultExtension(): String = "gox"
    override fun getIcon(): Icon? = GoxIcons.FILE
}