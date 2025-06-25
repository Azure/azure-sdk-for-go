import { TextDocument } from "vscode-languageserver-textdocument";
import { CompletionItem, Position } from "vscode-languageserver/node.js";
import { CompilerHost, ServerLog } from "../../index.js";
import { FileService } from "../file-service.js";
import { LibraryProvider } from "../lib-provider.js";
export declare function provideTspconfigCompletionItems(tspConfigDoc: TextDocument, tspConfigPosition: Position, context: {
    fileService: FileService;
    compilerHost: CompilerHost;
    emitterProvider: LibraryProvider;
    linterProvider: LibraryProvider;
    log: (log: ServerLog) => void;
}): Promise<CompletionItem[]>;
//# sourceMappingURL=completion.d.ts.map