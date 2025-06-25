import type { SystemHost } from "../core/types.js";
import type { InitTemplate, InitTemplateFile, InitTemplateLibrary, InitTemplateLibrarySpec } from "./init-template.js";
export declare const TypeSpecConfigFilename = "tspconfig.yaml";
export interface ScaffoldingConfig {
    /** Template used to resolve that config */
    template: InitTemplate;
    /**
     * Path where this template was loaded from.
     */
    baseUri: string;
    /**
     * Directory full path where the project should be initialized.
     */
    directory: string;
    /**
     * Name of the project.
     */
    name: string;
    /**
     * List of libraries to include
     */
    libraries: InitTemplateLibrarySpec[];
    /**
     * Whether to generate a .gitignore file.
     */
    includeGitignore: boolean;
    /**
     * Custom parameters provided in the templates.
     */
    parameters: Record<string, any>;
    /**
     * Selected emitters the tempalates.
     */
    emitters: Record<string, any>;
}
export declare function normalizeLibrary(library: InitTemplateLibrary): InitTemplateLibrarySpec;
export declare function makeScaffoldingConfig(template: InitTemplate, config: Partial<ScaffoldingConfig>): ScaffoldingConfig;
/**
 * Scaffold a new TypeSpec project using the given scaffolding config.
 * @param host
 * @param config
 */
export declare function scaffoldNewProject(host: SystemHost, config: ScaffoldingConfig): Promise<void>;
export declare function isFileSkipGeneration(fileName: string, files: InitTemplateFile[]): boolean;
//# sourceMappingURL=scaffold.d.ts.map