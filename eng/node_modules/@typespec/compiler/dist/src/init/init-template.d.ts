import type { JSONSchemaType } from "ajv";
import { TypeSpecRawConfig } from "../config/types.js";
export interface InitTemplateFile {
    path: string;
    destination: string;
    skipGeneration?: boolean;
}
export interface InitTemplateInput {
    description: string;
    type: "text";
    initialValue: any;
}
export interface InitTemplate {
    /**
     * The kind of project this tempolate initialize. This will change things like where dependencies are added.
     * For example, a library will add dependencies to `peer` and `dev` dependencies, while a project will add them to `dependencies`.
     * @default "project"
     */
    target?: "library" | "project";
    /**
     * Name of the template
     */
    title: string;
    /**
     * Description for the template.
     */
    description: string;
    /** Minimum Compiler Support Version */
    compilerVersion?: string;
    /**
     * List of libraries to include
     */
    libraries?: InitTemplateLibrary[];
    /**
     * List of emitters to include
     */
    emitters?: Record<string, EmitterTemplate>;
    /**
     * Config
     */
    config?: TypeSpecRawConfig;
    /**
     * Custom inputs to prompt to the user
     */
    inputs?: Record<string, InitTemplateInput>;
    /**
     * A flag to indicate not adding @typespec/compiler package to package.json.
     * Other libraries may already brought in the dependency such as Azure template.
     */
    skipCompilerPackage?: boolean;
    /**
     * List of files to copy.
     */
    files?: InitTemplateFile[];
}
/**
 * Describes emitter dependencies that will be added to the generated project.
 */
export interface EmitterTemplate {
    /** Friendly name for the emitter */
    label?: string;
    /** Emitter Selection Description */
    description?: string;
    /** Whether emitter is selected by default in the list */
    selected?: boolean;
    /** Optional emitter Options to populate the tspconfig.yaml */
    options?: any;
    /** Optional message to display to the user post creation */
    message?: string;
    /** Optional specific emitter version. `latest` if not specified */
    version?: string;
}
/**
 * Describes a library dependency that will be added to the generated project.
 */
export type InitTemplateLibrary = string | InitTemplateLibrarySpec;
/**
 * Describes a library dependency that will be added to the generated project.
 */
export interface InitTemplateLibrarySpec {
    /**
     * The npm package name of the library.
     */
    name: string;
    /**
     *  The npm-style version string as it would appear in package.json.
     */
    version?: string;
}
export declare const InitTemplateLibrarySpecSchema: JSONSchemaType<InitTemplateLibrarySpec>;
export declare const InitTemplateSchema: JSONSchemaType<InitTemplate>;
//# sourceMappingURL=init-template.d.ts.map