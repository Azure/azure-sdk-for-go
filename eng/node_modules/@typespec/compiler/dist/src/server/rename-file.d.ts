import { TextEdit } from "vscode-languageserver";
import { ImportStatementNode } from "../core/types.js";
/**
 * Get the LSP TextEdit to update an import statment value.
 * @param node Import statement node that should be updated
 * @param newImport New import path
 * @returns Lsp TextEdit
 */
export declare function getRenameImportEdit(node: ImportStatementNode, newImport: string): TextEdit;
export interface RenameFileParams {
    oldPath: string;
    newPath: string;
    isDirRename: boolean;
}
export interface ReplaceImportResult {
    newValue: string;
    filePath: string;
}
/**
 * Get the updated import value for a given rename operation.
 * If the rename operation is not applicable, it returns undefined.
 * @param target Import statement node that should be checked
 * @param params Current rename operation parameters
 * @returns Updated import value or undefined
 */
export declare function getUpdatedImportValue(target: ImportStatementNode, params: RenameFileParams): ReplaceImportResult | undefined;
//# sourceMappingURL=rename-file.d.ts.map