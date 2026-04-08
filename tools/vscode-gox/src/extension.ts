import * as vscode from 'vscode';
import * as path from 'path';
import * as fs from 'fs';
import { exec } from 'child_process';

let outputChannel: vscode.OutputChannel;

export function activate(context: vscode.ExtensionContext) {
    outputChannel = vscode.window.createOutputChannel('GoX');

    // Register commands
    context.subscriptions.push(
        vscode.commands.registerCommand('gox.compile', compileCurrentFile),
        vscode.commands.registerCommand('gox.compileAll', compileAllFiles),
        vscode.commands.registerCommand('gox.watch', watchAndCompile)
    );

    // Auto-compile on save
    const config = vscode.workspace.getConfiguration('gox');
    if (config.get<boolean>('compileOnSave')) {
        context.subscriptions.push(
            vscode.workspace.onDidSaveTextDocument(doc => {
                if (doc.fileName.endsWith('.gox')) {
                    compileFile(doc.fileName);
                }
            })
        );
    }

    outputChannel.appendLine('GoX extension activated');
}

function compileCurrentFile() {
    const editor = vscode.window.activeTextEditor;
    if (!editor || !editor.document.fileName.endsWith('.gox')) {
        vscode.window.showWarningMessage('No .gox file is currently open');
        return;
    }
    compileFile(editor.document.fileName);
}

function compileAllFiles() {
    const workspaceFolders = vscode.workspace.workspaceFolders;
    if (!workspaceFolders) {
        vscode.window.showWarningMessage('No workspace folder open');
        return;
    }

    const config = vscode.workspace.getConfiguration('gox');
    const goxPath = config.get<string>('goxPath') || 'gox';

    workspaceFolders.forEach(folder => {
        const cmd = `${goxPath} "${folder.uri.fsPath}"`;
        outputChannel.appendLine(`Running: ${cmd}`);

        exec(cmd, (error, stdout, stderr) => {
            if (error) {
                outputChannel.appendLine(`Error: ${error.message}`);
                vscode.window.showErrorMessage(`GoX compilation failed: ${error.message}`);
                return;
            }
            if (stderr) {
                outputChannel.appendLine(`Stderr: ${stderr}`);
            }
            if (stdout) {
                outputChannel.appendLine(`Stdout: ${stdout}`);
            }
            vscode.window.showInformationMessage('GoX compilation complete');
        });
    });
}

let watchProcess: any = null;

function watchAndCompile() {
    if (watchProcess) {
        watchProcess.kill();
        watchProcess = null;
        vscode.window.showInformationMessage('GoX watch stopped');
        return;
    }

    const workspaceFolders = vscode.workspace.workspaceFolders;
    if (!workspaceFolders) {
        vscode.window.showWarningMessage('No workspace folder open');
        return;
    }

    const config = vscode.workspace.getConfiguration('gox');
    const goxPath = config.get<string>('goxPath') || 'gox';

    workspaceFolders.forEach(folder => {
        const cmd = `${goxPath} -watch "${folder.uri.fsPath}"`;
        outputChannel.appendLine(`Running: ${cmd}`);

        const { spawn } = require('child_process');
        watchProcess = spawn(goxPath, ['-watch', folder.uri.fsPath]);

        watchProcess.stdout.on('data', (data: any) => {
            outputChannel.appendLine(`Stdout: ${data}`);
        });

        watchProcess.stderr.on('data', (data: any) => {
            outputChannel.appendLine(`Stderr: ${data}`);
        });

        watchProcess.on('close', (code: number) => {
            outputChannel.appendLine(`Process exited with code ${code}`);
            watchProcess = null;
        });
    });

    vscode.window.showInformationMessage('GoX watch started');
}

function compileFile(filePath: string) {
    const config = vscode.workspace.getConfiguration('gox');
    const goxPath = config.get<string>('goxPath') || 'gox';

    // Output path: .gox -> .go
    const outputPath = filePath.slice(0, -1) + 'o';

    const cmd = `${goxPath} -o "${outputPath}" "${filePath}"`;
    outputChannel.appendLine(`Running: ${cmd}`);

    exec(cmd, (error, stdout, stderr) => {
        if (error) {
            outputChannel.appendLine(`Error: ${error.message}`);
            vscode.window.showErrorMessage(`GoX compilation failed: ${error.message}`);
            return;
        }
        if (stderr) {
            outputChannel.appendLine(`Stderr: ${stderr}`);
        }
        if (stdout) {
            outputChannel.appendLine(`Stdout: ${stdout}`);
        }
        vscode.window.showInformationMessage(`Compiled: ${path.basename(outputPath)}`);
    });
}

export function deactivate() {
    if (watchProcess) {
        watchProcess.kill();
    }
    outputChannel.dispose();
}