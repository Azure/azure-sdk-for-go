import type { CompilerHost, SystemHost } from "../core/types.js";
export declare function mkTempDir(host: CompilerHost, base: string, prefix: string): Promise<string>;
/**
 * List all files in dir recursively
 * @returns relative path of the files from the given directory
 */
export declare function listAllFilesInDir(host: SystemHost, dir: string): Promise<string[]>;
//# sourceMappingURL=fs-utils.d.ts.map