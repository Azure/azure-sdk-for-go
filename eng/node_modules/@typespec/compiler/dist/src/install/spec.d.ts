import { DiagnosticError } from "../core/diagnostic-error.js";
import type { CompilerHost, SourceFile, Tracer } from "../core/types.js";
import { type SupportedPackageManager } from "./config.js";
export interface Descriptor {
    /** Name of the package manager */
    readonly name: SupportedPackageManager;
    /** Supported version range  */
    readonly range: string;
    readonly hash?: Hash;
}
export interface Hash {
    readonly algorithm: string;
    readonly value: string;
}
export interface ResolvedSpecResult {
    readonly kind: "resolved";
    /** Path to the resolved package.json */
    readonly path: string;
    /** Resolved spec of the package manager */
    readonly spec: Descriptor;
}
type PackageManagerSpecResult = {
    readonly kind: "no-package";
    readonly path: string;
} | {
    readonly kind: "no-spec";
    readonly path: string;
} | ResolvedSpecResult;
/**
 * Resolve the package manager required for the current working directory.
 * @throws {PackageManagerSpecError} if there is error resolving it(Invalid package.json, invalid packageManager field)
 */
export declare function resolvePackageManagerSpec(host: CompilerHost, parentTracer: Tracer, cwd: string): Promise<PackageManagerSpecResult>;
export declare function updatePackageManagerInPackageJson(host: CompilerHost, path: string, spec: Descriptor): Promise<void>;
export declare class PackageManagerSpecError extends DiagnosticError {
    constructor(message: string, source: SourceFile);
}
export {};
//# sourceMappingURL=spec.d.ts.map