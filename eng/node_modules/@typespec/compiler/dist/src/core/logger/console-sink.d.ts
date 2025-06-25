import { Diagnostic } from "../types.js";
export interface FormatDiagnosticOptions {
    readonly pretty?: boolean;
    readonly pathRelativeTo?: string;
}
export declare function formatDiagnostic(diagnostic: Diagnostic, options?: FormatDiagnosticOptions): string;
//# sourceMappingURL=console-sink.d.ts.map