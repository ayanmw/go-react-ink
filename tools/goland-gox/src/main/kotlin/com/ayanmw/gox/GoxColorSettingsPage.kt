package com.ayanmw.gox

import com.intellij.openapi.editor.colors.TextAttributesKey
import com.intellij.openapi.options.colors.AttributesDescriptor
import com.intellij.openapi.options.colors.ColorDescriptor
import com.intellij.openapi.options.colors.ColorSettingsPage
import com.intellij.openapi.util.NlsContexts
import javax.swing.Icon

/**
 * Color settings page for GoX
 */
class GoxColorSettingsPage : ColorSettingsPage {
    override fun getIcon(): Icon? = GoxIcons.FILE
    override fun getHighlighter() = GoxSyntaxHighlighter()
    override fun getDemoText(): String = """
// Package main demonstrates GoX syntax
package main

import (
    "github.com/ayanmw/go-react-ink/pkg/core"
    "github.com/ayanmw/go-react-ink/pkg/hooks"
)

// App is the main application
func App() core.Element {
    count, setCount := hooks.UseState(0)

    return <Box flexDirection="column" padding={1}>
        <Text color="green" bold>Count: {count}</Text>
        <Text dim>Press + to increment</Text>
        <Spacer />
        <Text>Footer</Text>
    </Box>
}

func main() {
    println("Hello, GoX!")
}
""".trimIndent()

    override fun getAdditionalHighlightingTagToDescriptorMap(): Map<String, TextAttributesKey>? = null
    override fun getAttributeDescriptors(): Array<AttributesDescriptor> = DESCRIPTORS
    override fun getColorDescriptors(): Array<ColorDescriptor> = ColorDescriptor.EMPTY_ARRAY
    override fun getDisplayName(): @NlsContexts.ConfigurableName String = "GoX"

    companion object {
        private val DESCRIPTORS = arrayOf(
            AttributesDescriptor("Keyword", GoxSyntaxHighlighter.KEYWORD),
            AttributesDescriptor("String", GoxSyntaxHighlighter.STRING),
            AttributesDescriptor("Number", GoxSyntaxHighlighter.NUMBER),
            AttributesDescriptor("Line comment", GoxSyntaxHighlighter.LINE_COMMENT),
            AttributesDescriptor("Block comment", GoxSyntaxHighlighter.BLOCK_COMMENT),
            AttributesDescriptor("Operator", GoxSyntaxHighlighter.OPERATOR),
            AttributesDescriptor("Brackets", GoxSyntaxHighlighter.BRACKET),
            AttributesDescriptor("Parentheses", GoxSyntaxHighlighter.PAREN),
            AttributesDescriptor("Braces", GoxSyntaxHighlighter.BRACE),
            AttributesDescriptor("JSX Tag", GoxSyntaxHighlighter.JSX_TAG),
            AttributesDescriptor("JSX Tag Name", GoxSyntaxHighlighter.JSX_TAG_NAME),
            AttributesDescriptor("JSX Attribute Name", GoxSyntaxHighlighter.JSX_ATTR_NAME),
            AttributesDescriptor("JSX Attribute Value", GoxSyntaxHighlighter.JSX_ATTR_VALUE),
            AttributesDescriptor("JSX Expression", GoxSyntaxHighlighter.JSX_EXPR),
            AttributesDescriptor("Identifier", GoxSyntaxHighlighter.IDENTIFIER),
            AttributesDescriptor("Function name", GoxSyntaxHighlighter.FUNCTION_NAME)
        )
    }
}