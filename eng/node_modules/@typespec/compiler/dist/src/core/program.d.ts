import { Checker } from "./checker.js";
import { CompilerOptions } from "./options.js";
import { CompilerHost, Diagnostic, EmitterFunc, JsSourceFileNode, LibraryInstance, LibraryMetadata, LocationContext, Namespace, SourceFile, Tracer, Type, TypeSpecScriptNode } from "./types.js";
export interface Program {
    compilerOptions: CompilerOptions;
    /** All source files in the program, keyed by their file path. */
    sourceFiles: Map<string, TypeSpecScriptNode>;
    jsSourceFiles: Map<string, JsSourceFileNode>;
    host: CompilerHost;
    tracer: Tracer;
    trace(area: string, message: string): void;
    /**
     * **DANGER** Using the checker is reserved for advanced usage and should be used with caution.
     * API are not subject to the same stability guarantees see See https://typespec.io/docs/handbook/breaking-change-policy/
     */
    checker: Checker;
    emitters: EmitterRef[];
    readonly diagnostics: readonly Diagnostic[];
    stateSet(key: symbol): Set<Type>;
    stateMap(key: symbol): Map<Type, any>;
    hasError(): boolean;
    reportDiagnostic(diagnostic: Diagnostic): void;
    reportDiagnostics(diagnostics: readonly Diagnostic[]): void;
    getGlobalNamespaceType(): Namespace;
    resolveTypeReference(reference: string): [Type | undefined, readonly Diagnostic[]];
    /** Return location context of the given source file. */
    getSourceFileLocationContext(sourceFile: SourceFile): LocationContext;
    /**
     * Project root. If a tsconfig was found/specified this is the directory for the tsconfig.json. Otherwise directory where the entrypoint is located.
     */
    readonly projectRoot: string;
}
interface EmitterRef {
    emitFunction: EmitterFunc;
    main: string;
    metadata: LibraryMetadata;
    emitterOutputDir: string;
    options: Record<string, unknown>;
    readonly library: LibraryInstance;
}
export declare function compile(host: CompilerHost, mainFile: string, options?: CompilerOptions, oldProgram?: Program): Promise<Program>;
export {};
//# sourceMappingURL=program.d.ts.map