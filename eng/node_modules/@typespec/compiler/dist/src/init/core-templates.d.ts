import type { SystemHost } from "../core/types.js";
export declare const templatesDir: string;
export interface LoadedCoreTemplates {
    readonly baseUri: string;
    readonly templates: Record<string, any>;
}
export declare function getTypeSpecCoreTemplates(host: SystemHost): Promise<LoadedCoreTemplates>;
//# sourceMappingURL=core-templates.d.ts.map