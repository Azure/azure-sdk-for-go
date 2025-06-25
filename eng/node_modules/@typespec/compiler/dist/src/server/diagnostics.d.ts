import type { Diagnostic as VSDiagnostic } from "vscode-languageserver";
import type { TextDocument } from "vscode-languageserver-textdocument";
import type { Program } from "../core/program.js";
import { Diagnostic } from "../core/types.js";
import type { FileService } from "./file-service.js";
/** Convert TypeSpec Diagnostic to Lsp diagnostic. Each TypeSpec diagnostic could produce multiple lsp ones when it involve multiple locations. */
export declare function convertDiagnosticToLsp(fileService: FileService, program: Program, document: TextDocument, diagnostic: Diagnostic): [VSDiagnostic, TextDocument][];
//# sourceMappingURL=diagnostics.d.ts.map