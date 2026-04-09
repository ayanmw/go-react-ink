"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (k !== "default" && Object.prototype.hasOwnProperty.call(mod, k)) __createBinding(result, mod, k);
    __setModuleDefault(result, mod);
    return result;
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.deactivate = exports.activate = void 0;
const vscode = __importStar(require("vscode"));
const path = __importStar(require("path"));
const child_process_1 = require("child_process");
let outputChannel;
let lspClient = null;
function activate(context) {
    outputChannel = vscode.window.createOutputChannel('GoX');
    // Register commands
    context.subscriptions.push(vscode.commands.registerCommand('gox.compile', compileCurrentFile), vscode.commands.registerCommand('gox.compileAll', compileAllFiles), vscode.commands.registerCommand('gox.watch', watchAndCompile), vscode.commands.registerCommand('gox.lint', runLint));
    // Start LSP client if enabled
    const config = vscode.workspace.getConfiguration('gox');
    if (config.get('lsp.enabled')) {
        startLSPClient(context);
    }
    // Auto-compile on save
    if (config.get('compileOnSave')) {
        context.subscriptions.push(vscode.workspace.onDidSaveTextDocument(doc => {
            if (doc.fileName.endsWith('.gox')) {
                compileFile(doc.fileName);
            }
        }));
    }
    outputChannel.appendLine('GoX extension activated');
}
exports.activate = activate;
function startLSPClient(context) {
    const config = vscode.workspace.getConfiguration('gox');
    const lspPath = config.get('lsp.path') || 'gox-lsp';
    lspClient = new LSPClient(lspPath, outputChannel);
    lspClient.start();
    context.subscriptions.push({
        dispose: () => {
            if (lspClient) {
                lspClient.stop();
            }
        }
    });
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
    const goxPath = config.get('goxPath') || 'gox';
    workspaceFolders.forEach(folder => {
        const cmd = `${goxPath} "${folder.uri.fsPath}"`;
        outputChannel.appendLine(`Running: ${cmd}`);
        (0, child_process_1.exec)(cmd, (error, stdout, stderr) => {
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
let watchProcess = null;
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
    const goxPath = config.get('goxPath') || 'gox';
    workspaceFolders.forEach(folder => {
        const cmd = `${goxPath} -watch "${folder.uri.fsPath}"`;
        outputChannel.appendLine(`Running: ${cmd}`);
        const { spawn } = require('child_process');
        watchProcess = spawn(goxPath, ['-watch', folder.uri.fsPath]);
        watchProcess.stdout.on('data', (data) => {
            outputChannel.appendLine(`Stdout: ${data}`);
        });
        watchProcess.stderr.on('data', (data) => {
            outputChannel.appendLine(`Stderr: ${data}`);
        });
        watchProcess.on('close', (code) => {
            outputChannel.appendLine(`Process exited with code ${code}`);
            watchProcess = null;
        });
    });
    vscode.window.showInformationMessage('GoX watch started');
}
function runLint() {
    const workspaceFolders = vscode.workspace.workspaceFolders;
    if (!workspaceFolders) {
        vscode.window.showWarningMessage('No workspace folder open');
        return;
    }
    outputChannel.appendLine('=== Running Lint ===');
    outputChannel.show();
    const workspacePath = workspaceFolders[0].uri.fsPath;
    let hasErrors = false;
    // Run gofmt
    outputChannel.appendLine('Checking gofmt...');
    (0, child_process_1.exec)('gofmt -s -d .', { cwd: workspacePath }, (error, stdout, stderr) => {
        if (stdout) {
            outputChannel.appendLine('gofmt found issues:');
            outputChannel.appendLine(stdout);
            hasErrors = true;
        }
        else {
            outputChannel.appendLine('✓ gofmt: No formatting issues');
        }
        // Run go vet
        outputChannel.appendLine('');
        outputChannel.appendLine('Running go vet...');
        (0, child_process_1.exec)('go vet ./...', { cwd: workspacePath }, (error, stdout, stderr) => {
            if (error) {
                outputChannel.appendLine('go vet found issues:');
                outputChannel.appendLine(stderr || stdout);
                hasErrors = true;
            }
            else {
                outputChannel.appendLine('✓ go vet: No issues');
            }
            // Run golint
            outputChannel.appendLine('');
            outputChannel.appendLine('Running golint...');
            (0, child_process_1.exec)('golint ./...', { cwd: workspacePath }, (error, stdout, stderr) => {
                if (stdout) {
                    outputChannel.appendLine('golint warnings:');
                    outputChannel.appendLine(stdout);
                }
                else {
                    outputChannel.appendLine('✓ golint: No issues');
                }
                outputChannel.appendLine('');
                if (hasErrors) {
                    outputChannel.appendLine('=== Lint completed with errors ===');
                    vscode.window.showErrorMessage('Lint found issues. Check output for details.');
                }
                else {
                    outputChannel.appendLine('=== Lint completed successfully ===');
                    vscode.window.showInformationMessage('Lint passed!');
                }
            });
        });
    });
}
function compileFile(filePath) {
    const config = vscode.workspace.getConfiguration('gox');
    const goxPath = config.get('goxPath') || 'gox';
    // Output path: .gox -> .go
    const outputPath = filePath.slice(0, -1) + 'o';
    const cmd = `${goxPath} -o "${outputPath}" "${filePath}"`;
    outputChannel.appendLine(`Running: ${cmd}`);
    (0, child_process_1.exec)(cmd, (error, stdout, stderr) => {
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
// Simple LSP Client implementation
class LSPClient {
    constructor(lspPath, outputChannel) {
        this.process = null;
        this.pendingRequests = new Map();
        this.requestId = 0;
        this.buffer = '';
        this.lspPath = lspPath;
        this.outputChannel = outputChannel;
    }
    start() {
        const { spawn } = require('child_process');
        this.process = spawn(this.lspPath, [], {
            stdio: ['pipe', 'pipe', 'pipe']
        });
        this.process.stdout.on('data', (data) => {
            this.handleData(data.toString());
        });
        this.process.stderr.on('data', (data) => {
            this.outputChannel.appendLine(`LSP stderr: ${data}`);
        });
        this.process.on('close', (code) => {
            this.outputChannel.appendLine(`LSP process exited with code ${code}`);
            this.process = null;
        });
        // Send initialize request
        this.sendRequest('initialize', {
            processId: process.pid,
            rootUri: vscode.workspace.workspaceFolders?.[0]?.uri.toString(),
            capabilities: {
                textDocument: {
                    completion: {
                        completionItem: {
                            snippetSupport: true
                        }
                    },
                    hover: {
                        contentFormat: ['markdown', 'plaintext']
                    }
                }
            }
        });
        this.outputChannel.appendLine('LSP client started');
    }
    stop() {
        if (this.process) {
            this.sendRequest('shutdown', null);
            this.process.kill();
            this.process = null;
        }
    }
    handleData(data) {
        this.buffer += data;
        while (true) {
            const headerEnd = this.buffer.indexOf('\r\n\r\n');
            if (headerEnd === -1)
                break;
            const header = this.buffer.substring(0, headerEnd);
            const contentLengthMatch = header.match(/Content-Length: (\d+)/);
            if (!contentLengthMatch)
                break;
            const contentLength = parseInt(contentLengthMatch[1]);
            const contentStart = headerEnd + 4;
            const contentEnd = contentStart + contentLength;
            if (this.buffer.length < contentEnd)
                break;
            const content = this.buffer.substring(contentStart, contentEnd);
            this.buffer = this.buffer.substring(contentEnd);
            this.handleMessage(JSON.parse(content));
        }
    }
    handleMessage(message) {
        if (message.method === 'textDocument/publishDiagnostics') {
            const uri = message.params.uri;
            const diagnostics = message.params.diagnostics;
            const doc = vscode.workspace.textDocuments.find(d => d.uri.toString() === uri);
            if (doc) {
                const diagnosticCollection = vscode.languages.createDiagnosticCollection('gox');
                diagnosticCollection.set(doc.uri, diagnostics.map((d) => {
                    const range = new vscode.Range(new vscode.Position(d.range.start.line, d.range.start.character), new vscode.Position(d.range.end.line, d.range.end.character));
                    const severity = d.severity === 1 ? vscode.DiagnosticSeverity.Error :
                        d.severity === 2 ? vscode.DiagnosticSeverity.Warning :
                            vscode.DiagnosticSeverity.Information;
                    return new vscode.Diagnostic(range, d.message, severity);
                }));
            }
        }
    }
    sendRequest(method, params) {
        const content = JSON.stringify({
            jsonrpc: '2.0',
            id: this.requestId++,
            method: method,
            params: params
        });
        const header = `Content-Length: ${content.length}\r\n\r\n`;
        this.process.stdin.write(header + content);
    }
    sendNotification(method, params) {
        const content = JSON.stringify({
            jsonrpc: '2.0',
            method: method,
            params: params
        });
        const header = `Content-Length: ${content.length}\r\n\r\n`;
        this.process.stdin.write(header + content);
    }
    // Notify LSP of document open
    documentOpened(doc) {
        this.sendNotification('textDocument/didOpen', {
            textDocument: {
                uri: doc.uri.toString(),
                languageId: 'gox',
                version: doc.version,
                text: doc.getText()
            }
        });
    }
    // Notify LSP of document change
    documentChanged(doc, content) {
        this.sendNotification('textDocument/didChange', {
            textDocument: {
                uri: doc.uri.toString(),
                version: doc.version
            },
            contentChanges: [{ text: content }]
        });
    }
    // Notify LSP of document close
    documentClosed(doc) {
        this.sendNotification('textDocument/didClose', {
            textDocument: {
                uri: doc.uri.toString()
            }
        });
    }
}
function deactivate() {
    if (watchProcess) {
        watchProcess.kill();
    }
    if (lspClient) {
        lspClient.stop();
    }
    outputChannel.dispose();
}
exports.deactivate = deactivate;
//# sourceMappingURL=extension.js.map