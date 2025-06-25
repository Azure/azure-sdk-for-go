import type { CliCompilerHost } from "../core/cli/types.js";
import { DiagnosticError } from "../core/diagnostic-error.js";
import { type Diagnostic } from "../core/types.js";
export declare class InstallDependenciesError extends DiagnosticError {
    constructor(message: string);
}
export interface InstallTypeSpecDependenciesOptions {
    readonly directory: string;
    readonly stdio?: "inherit" | "pipe";
    /** When set to true update the packageManager field with the package manger version and hash */
    readonly savePackageManager?: boolean;
}
export declare function installTypeSpecDependencies(host: CliCompilerHost, options: InstallTypeSpecDependenciesOptions): Promise<readonly Diagnostic[]>;
//# sourceMappingURL=install.d.ts.map