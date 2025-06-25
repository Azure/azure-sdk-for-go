import type { Diagnostic } from "./types.js";
/** Error that reuse the diagnostic system. */
export declare class DiagnosticError extends Error {
    readonly diagnostics: readonly Diagnostic[];
    constructor(diagnostics: Diagnostic | readonly Diagnostic[]);
}
//# sourceMappingURL=diagnostic-error.d.ts.map