package com.ayanmw.gox.actions

import com.intellij.execution.configurations.GeneralCommandLine
import com.intellij.execution.process.CapturingProcessHandler
import com.intellij.notification.Notification
import com.intellij.notification.NotificationType
import com.intellij.notification.Notifications
import com.intellij.openapi.actionSystem.AnAction
import com.intellij.openapi.actionSystem.AnActionEvent
import com.intellij.openapi.actionSystem.CommonDataKeys
import com.intellij.openapi.progress.ProgressIndicator
import com.intellij.openapi.progress.ProgressManager
import com.intellij.openapi.progress.Task
import com.intellij.openapi.project.Project
import com.intellij.openapi.vfs.VirtualFile
import com.intellij.psi.PsiFile
import java.io.File

/**
 * Action to compile current .gox file
 */
class CompileAction : AnAction() {
    override fun actionPerformed(e: AnActionEvent) {
        val project = e.project ?: return
        val file = e.getData(CommonDataKeys.PSI_FILE) ?: return

        if (file.virtualFile.extension != "gox") {
            showNotification(project, "Not a .gox file", NotificationType.WARNING)
            return
        }

        ProgressManager.getInstance().run(object : Task.Backgroundable(project, "Compiling .gox file") {
            override fun run(indicator: ProgressIndicator) {
                compileFile(project, file.virtualFile)
            }
        })
    }

    override fun update(e: AnActionEvent) {
        val file = e.getData(CommonDataKeys.PSI_FILE)
        e.presentation.isEnabled = file?.virtualFile?.extension == "gox"
    }

    private fun compileFile(project: Project, file: VirtualFile) {
        val inputPath = file.path
        val outputPath = inputPath.removeSuffix(".gox") + ".go"

        try {
            val commandLine = GeneralCommandLine("gox", "-o", outputPath, inputPath)
            commandLine.setWorkDirectory(project.basePath)

            val handler = CapturingProcessHandler(commandLine)
            val output = handler.runProcess(30, true)

            if (output.exitCode == 0) {
                showNotification(project, "Compiled: ${file.name} → ${File(outputPath).name}", NotificationType.INFORMATION)
            } else {
                showNotification(project, "Compilation failed: ${output.stderr}", NotificationType.ERROR)
            }
        } catch (e: Exception) {
            showNotification(project, "Error: ${e.message}", NotificationType.ERROR)
        }
    }

    private fun showNotification(project: Project, message: String, type: NotificationType) {
        val notification = Notification("GoX", "GoX Compiler", message, type)
        Notifications.Bus.notify(notification, project)
    }
}

/**
 * Action to compile all .gox files in project
 */
class CompileAllAction : AnAction() {
    override fun actionPerformed(e: AnActionEvent) {
        val project = e.project ?: return

        ProgressManager.getInstance().run(object : Task.Backgroundable(project, "Compiling all .gox files") {
            override fun run(indicator: ProgressIndicator) {
                compileAll(project, indicator)
            }
        })
    }

    private fun compileAll(project: Project, indicator: ProgressIndicator) {
        val basePath = project.basePath ?: return

        try {
            val commandLine = GeneralCommandLine("gox", basePath)
            commandLine.setWorkDirectory(basePath)

            val handler = CapturingProcessHandler(commandLine)
            val output = handler.runProcess(60, true)

            if (output.exitCode == 0) {
                showNotification(project, "All .gox files compiled", NotificationType.INFORMATION)
            } else {
                showNotification(project, "Compilation failed: ${output.stderr}", NotificationType.ERROR)
            }
        } catch (e: Exception) {
            showNotification(project, "Error: ${e.message}", NotificationType.ERROR)
        }
    }

    private fun showNotification(project: Project, message: String, type: NotificationType) {
        val notification = Notification("GoX", "GoX Compiler", message, type)
        Notifications.Bus.notify(notification, project)
    }
}