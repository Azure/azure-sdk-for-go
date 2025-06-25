import { ignoreDiagnostics } from "../core/diagnostics.js";
/**
 * Creates a diagnosable function wrapper.
 *
 * The returned function will ignore diagnostics by default.
 * A `withDiagnostics` property is attached to the returned function,
 * allowing access to the original function's return value including diagnostics.
 *
 * @param fn The function to wrap, which must return a tuple `[Result, readonly Diagnostic[]]`.
 * @returns A function that ignores diagnostics by default, with a `withDiagnostics` method.
 */
export function createDiagnosable(fn) {
    function wrapper(...args) {
        return ignoreDiagnostics(fn.apply(this, args));
    }
    wrapper.withDiagnostics = fn;
    return wrapper;
}
//# sourceMappingURL=create-diagnosable.js.map