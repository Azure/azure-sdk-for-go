import { mkdir, stat } from "fs/promises";
import { fileURLToPath, pathToFileURL } from "url";
import { findProjectRoot } from "../utils/io.js";
import { createConsoleSink } from "./logger/index.js";
import { NodeSystemHost } from "./node-system-host.js";
import { joinPaths } from "./path-utils.js";
import { getSourceFileKindFromExt } from "./source-file.js";
export const CompilerPackageRoot = (await findProjectRoot(stat, fileURLToPath(import.meta.url)));
/**
 * Implementation of the @see CompilerHost using the real file system.
 * This is the the CompilerHost used by TypeSpec CLI.
 */
export const NodeHost = {
    ...NodeSystemHost,
    getExecutionRoot: () => CompilerPackageRoot,
    getJsImport: (path) => import(pathToFileURL(path).href),
    getLibDirs() {
        const rootDir = this.getExecutionRoot();
        return [joinPaths(rootDir, "lib/std")];
    },
    getSourceFileKind: getSourceFileKindFromExt,
    mkdirp: (path) => mkdir(path, { recursive: true }),
    logSink: createConsoleSink(),
    fileURLToPath,
    pathToFileURL(path) {
        return pathToFileURL(path).href;
    },
};
//# sourceMappingURL=node-host.js.map