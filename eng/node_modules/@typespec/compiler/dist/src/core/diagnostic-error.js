/** Error that reuse the diagnostic system. */
export class DiagnosticError extends Error {
    diagnostics;
    constructor(diagnostics) {
        super();
        this.diagnostics = Array.isArray(diagnostics) ? diagnostics : [diagnostics];
    }
}
//# sourceMappingURL=diagnostic-error.js.map