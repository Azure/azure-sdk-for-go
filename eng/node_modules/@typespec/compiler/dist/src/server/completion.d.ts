import { CompletionList, CompletionParams } from "vscode-languageserver";
import { Program } from "../core/program.js";
import { PositionDetail, TypeSpecScriptNode } from "../core/types.js";
export type CompletionContext = {
    program: Program;
    params: CompletionParams;
    file: TypeSpecScriptNode;
    completions: CompletionList;
};
export declare function resolveCompletion(context: CompletionContext, posDetail: PositionDetail): Promise<CompletionList>;
//# sourceMappingURL=completion.d.ts.map