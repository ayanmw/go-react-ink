package com.ayanmw.gox

import com.intellij.codeInsight.completion.*
import com.intellij.codeInsight.lookup.LookupElementBuilder
import com.intellij.patterns.PlatformPatterns
import com.intellij.util.ProcessingContext

/**
 * Completion contributor for GoX
 */
class GoxCompletionContributor : CompletionContributor() {
    init {
        // JSX Component completion
        extend(CompletionType.BASIC,
            PlatformPatterns.psiElement().withLanguage(GoxLanguage),
            object : CompletionProvider<CompletionParameters>() {
                override fun addCompletions(
                    parameters: CompletionParameters,
                    context: ProcessingContext,
                    result: CompletionResultSet
                ) {
                    val position = parameters.position
                    val text = position.containingFile.text
                    val offset = position.textOffset

                    // Check if we're inside JSX context
                    val inJSX = isInJSXContext(text, offset)

                    if (inJSX) {
                        // Add component completions
                        COMPONENTS.forEach { (name, desc) ->
                            result.addElement(
                                LookupElementBuilder.create(name)
                                    .withTypeText(desc)
                                    .withInsertHandler { context, _ ->
                                        val editor = context.editor
                                        val document = editor.document
                                        val offset = context.tailOffset
                                        document.insertString(offset, ">$1</$name>")
                                        editor.caretModel.moveToOffset(offset + 1)
                                    }
                            )
                        }
                    }

                    // Add hook completions
                    HOOKS.forEach { (name, desc) ->
                        result.addElement(
                            LookupElementBuilder.create(name)
                                .withTypeText(desc)
                                .withIcon(com.intellij.icons.AllIcons.Nodes.Function)
                        )
                    }

                    // Add attribute completions for JSX
                    if (inJSX && isAfterTagName(text, offset)) {
                        ATTRIBUTES.forEach { (name, desc) ->
                            result.addElement(
                                LookupElementBuilder.create(name)
                                    .withTypeText(desc)
                                    .withInsertHandler { context, _ ->
                                        val editor = context.editor
                                        val document = editor.document
                                        val offset = context.tailOffset
                                        document.insertString(offset, "={}")
                                        editor.caretModel.moveToOffset(offset + 2)
                                    }
                            )
                        }
                    }
                }
            }
        )
    }

    private fun isInJSXContext(text: String, offset: Int): Boolean {
        var depth = 0
        var i = 0
        while (i < offset) {
            when (text[i]) {
                '<' -> if (i + 1 < text.length && text[i + 1].isUpperCase()) depth++
                '>' -> if (depth > 0) depth--
            }
            i++
        }
        return depth > 0
    }

    private fun isAfterTagName(text: String, offset: Int): Boolean {
        var i = offset - 1
        while (i >= 0 && text[i].isWhitespace()) i--
        return i >= 0 && (text[i].isLetterOrDigit() || text[i] == '>')
    }

    companion object {
        val COMPONENTS = mapOf(
            "Box" to "Container component",
            "Text" to "Text component",
            "Spacer" to "Flexible space",
            "Newline" to "Line break",
            "Static" to "Static content",
            "Transform" to "Transform wrapper",
            "Fragment" to "Fragment container"
        )

        val HOOKS = mapOf(
            "useState" to "State hook",
            "useEffect" to "Effect hook",
            "useLayoutEffect" to "Layout effect hook",
            "useRef" to "Reference hook",
            "useMemo" to "Memoization hook",
            "useCallback" to "Callback hook",
            "useInput" to "Input hook",
            "useApp" to "App context hook",
            "useFocus" to "Focus hook",
            "useFocusManager" to "Focus manager hook",
            "useCursor" to "Cursor hook",
            "useAnimation" to "Animation hook",
            "useStdout" to "Stdout hook",
            "useStdin" to "Stdin hook",
            "useStderr" to "Stderr hook",
            "useWindowSize" to "Window size hook",
            "useBoxMetrics" to "Box metrics hook",
            "usePaste" to "Paste hook",
            "useIsScreenReaderEnabled" to "Screen reader hook"
        )

        val ATTRIBUTES = mapOf(
            "flexDirection" to "row | column",
            "justifyContent" to "flex-start | center | flex-end",
            "alignItems" to "flex-start | center | flex-end",
            "width" to "Width",
            "height" to "Height",
            "padding" to "Padding",
            "margin" to "Margin",
            "color" to "Text color",
            "backgroundColor" to "Background color",
            "bold" to "Bold text",
            "italic" to "Italic text",
            "underline" to "Underline text"
        )
    }
}