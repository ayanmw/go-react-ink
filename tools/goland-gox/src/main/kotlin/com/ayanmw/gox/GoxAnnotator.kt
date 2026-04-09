package com.ayanmw.gox

import com.intellij.lang.annotation.AnnotationHolder
import com.intellij.lang.annotation.Annotator
import com.intellij.lang.annotation.HighlightSeverity
import com.intellij.psi.PsiElement

/**
 * Annotator for GoX - provides error highlighting
 */
class GoxAnnotator : Annotator {
    override fun annotate(element: PsiElement, holder: AnnotationHolder) {
        // Basic validation - check for unclosed JSX tags
        if (element.text.startsWith("<") && !element.text.contains(">")) {
            holder.newAnnotation(HighlightSeverity.ERROR, "Unclosed JSX tag")
                .range(element.textRange)
                .create()
        }
    }
}