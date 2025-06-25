import { type Diagnostic } from "../core/types.js";
/**
 * Assert there is no diagnostics.
 * @param diagnostics Diagnostics
 */
export declare function expectDiagnosticEmpty(diagnostics: readonly Diagnostic[]): void;
/**
 * Condition to match
 */
export interface DiagnosticMatch {
    /**
     * Match the code.
     */
    code?: string;
    /**
     * Match the message.
     */
    message?: string | RegExp;
    /**
     * Match the severity.
     */
    severity?: "error" | "warning";
    /**
     * Name of the file for this diagnostic.
     */
    file?: string | RegExp;
    /**
     * Start position of the diagnostic
     */
    pos?: number;
    /**
     * End position of the diagnostic
     */
    end?: number;
}
/**
 * Validate the diagnostic array contains exactly the given diagnostics.
 * @param diagnostics Array of the diagnostics
 */
export declare function expectDiagnostics(diagnostics: readonly Diagnostic[], match: DiagnosticMatch | DiagnosticMatch[], options?: {
    strict: boolean;
}): void;
//# sourceMappingURL=expect.d.ts.map