import { Imports } from "../../types/package-json.js";
import { EsmResolutionContext } from "./utils.js";
/** Implementation of PACKAGE_IMPORTS_RESOLVE https://github.com/nodejs/node/blob/main/doc/api/esm.md */
export declare function resolvePackageImports(context: EsmResolutionContext, imports: Imports): Promise<string | null | undefined>;
//# sourceMappingURL=resolve-package-imports.d.ts.map