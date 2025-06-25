import { Diagnostic } from "./types.js";
export interface TypeSpecFormatOptions {
    exclude?: string[];
}
/**  The paths of which are either relative or absolute based on the original file path patterns. */
export interface TypeSpecFormatResult {
    /** Files which were formatted successfully, */
    readonly formatted: string[];
    /** Files which had a valid format already. */
    readonly alreadyFormatted: string[];
    /** Files that were included in the filter but are not in the scope of the typespec formatter. */
    readonly ignored: string[];
    /** Files with errors */
    readonly errored: [string, Diagnostic][];
}
/**
 * Format all the TypeSpec project files(.tsp, tspconfig.yaml).
 * @param patterns List of wildcard pattern searching for TypeSpec files.
 * @returns list of files which failed to format.
 */
export declare function formatFiles(patterns: string[], { exclude }: TypeSpecFormatOptions): Promise<TypeSpecFormatResult>;
export interface CheckFilesFormatResult {
    readonly formatted: string[];
    readonly needsFormat: string[];
    readonly ignored: string[];
    readonly errored: [string, Diagnostic][];
}
/**
 * Check the format of the files in the given pattern.
 */
export declare function checkFilesFormat(patterns: string[], { exclude }: TypeSpecFormatOptions): Promise<CheckFilesFormatResult>;
export type FormatFileResult = 
/** File formatted successfully. */
{
    kind: "formatted";
}
/** File was already formatted. */
 | {
    kind: "already-formatted";
}
/** File is not in a format that can be formatted by TypeSpec */
 | {
    kind: "ignored";
}
/** Error occurred, probably a parsing error. */
 | {
    kind: "error";
    diagnostic: Diagnostic;
};
export declare function formatFile(filename: string): Promise<FormatFileResult>;
export type CheckFormatResult = 
/** File formatted successfully. */
{
    kind: "formatted";
}
/** File needs format */
 | {
    kind: "needs-format";
}
/** File is not in a format that can be formatted by TypeSpec */
 | {
    kind: "ignored";
}
/** Error occurred, probably a parsing error. */
 | {
    kind: "error";
    diagnostic: Diagnostic;
};
/**
 * Check the given TypeSpec file is correctly formatted.
 */
export declare function checkFileFormat(filename: string): Promise<CheckFormatResult>;
//# sourceMappingURL=formatter-fs.d.ts.map