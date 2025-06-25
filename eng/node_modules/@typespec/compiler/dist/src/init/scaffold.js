import { stringify } from "yaml";
import { getDirectoryPath, joinPaths } from "../core/path-utils.js";
import { readUrlOrPath, resolveRelativeUrlOrPath } from "../utils/misc.js";
import { createFileTemplatingContext, render, } from "./file-templating.js";
export const TypeSpecConfigFilename = "tspconfig.yaml";
export function normalizeLibrary(library) {
    if (typeof library === "string") {
        return { name: library };
    }
    return library;
}
export function makeScaffoldingConfig(template, config) {
    return {
        template,
        libraries: config.libraries ?? template.libraries?.map(normalizeLibrary) ?? [],
        baseUri: config.baseUri ?? ".",
        name: config.name ?? "",
        directory: config.directory ?? "",
        parameters: config.parameters ?? {},
        includeGitignore: config.includeGitignore ?? true,
        emitters: config.emitters ?? {},
        ...config,
    };
}
/**
 * Scaffold a new TypeSpec project using the given scaffolding config.
 * @param host
 * @param config
 */
export async function scaffoldNewProject(host, config) {
    await host.mkdirp(config.directory);
    await writePackageJson(host, config);
    await writeConfig(host, config);
    await writeMain(host, config);
    await writeGitignore(host, config);
    await writeFiles(host, config);
}
export function isFileSkipGeneration(fileName, files) {
    for (const file of files) {
        if (file.destination === fileName) {
            return file.skipGeneration ?? false;
        }
    }
    return false;
}
async function writePackageJson(host, config) {
    if (isFileSkipGeneration("package.json", config.template.files ?? [])) {
        return;
    }
    const dependencies = {};
    if (!config.template.skipCompilerPackage) {
        dependencies["@typespec/compiler"] = "latest";
    }
    for (const library of config.libraries) {
        dependencies[library.name] = await getPackageVersion(library);
    }
    for (const key of Object.keys(config.emitters)) {
        dependencies[key] = await getPackageVersion(config.emitters[key]);
    }
    const packageJson = {
        name: config.name,
        version: "0.1.0",
        type: "module",
        private: true,
    };
    if (config.template.target === "library") {
        packageJson.peerDependencies = dependencies;
        packageJson.devDependencies = dependencies;
    }
    else {
        packageJson.dependencies = dependencies;
    }
    return host.writeFile(joinPaths(config.directory, "package.json"), JSON.stringify(packageJson, null, 2));
}
const placeholderConfig = `
# extends: ../tspconfig.yaml                    # Extend another config file
# emit:                                         # Emitter name
#   - "<emitter-name"
# options:                                      # Emitter options
#   <emitter-name>:
#    "<option-name>": "<option-value>"
# environment-variables:                        # Environment variables which can be used to interpolate emitter options
#   <variable-name>:
#     default: "<variable-default>"
# parameters:                                   # Parameters which can be used to interpolate emitter options
#   <param-name>:
#     default: "<param-default>"
# trace:                                        # Trace areas to enable tracing
#  - "<trace-name>"
# warn-as-error: true                           # Treat warnings as errors
# output-dir: "{project-root}/_generated"       # Configure the base output directory for all emitters
`.trim();
async function writeConfig(host, config) {
    if (isFileSkipGeneration(TypeSpecConfigFilename, config.template.files ?? [])) {
        return;
    }
    let rawConfig;
    if (config.template.config !== undefined && Object.keys(config.template.config).length > 0) {
        rawConfig = config.template.config;
    }
    if (Object.keys(config.emitters).length > 0) {
        rawConfig ??= {};
        rawConfig.emit = Object.keys(config.emitters);
        rawConfig.options = Object.fromEntries(Object.entries(config.emitters).map(([key, emitter]) => [key, emitter.options]));
    }
    const content = rawConfig ? stringify(rawConfig) : placeholderConfig;
    return host.writeFile(joinPaths(config.directory, TypeSpecConfigFilename), content);
}
async function writeMain(host, config) {
    if (isFileSkipGeneration("main.tsp", config.template.files ?? [])) {
        return;
    }
    const dependencies = {};
    for (const library of config.libraries) {
        dependencies[library.name] = await getPackageVersion(library);
    }
    const lines = [...config.libraries.map((x) => `import "${x.name}";`), ""];
    const content = lines.join("\n");
    return host.writeFile(joinPaths(config.directory, "main.tsp"), content);
}
const defaultGitignore = `
# MacOS
.DS_Store

# Default TypeSpec output
tsp-output/
dist/

# Dependency directories
node_modules/
`.trim();
async function writeGitignore(host, config) {
    if (!config.includeGitignore || isFileSkipGeneration(".gitignore", config.template.files ?? [])) {
        return;
    }
    return host.writeFile(joinPaths(config.directory, ".gitignore"), defaultGitignore);
}
async function writeFiles(host, config) {
    const templateContext = createFileTemplatingContext(config);
    if (!config.template.files) {
        return;
    }
    for (const file of config.template.files) {
        if (file.skipGeneration !== true) {
            await writeFile(host, config, templateContext, file);
        }
    }
}
async function writeFile(host, config, context, file) {
    const baseDir = config.baseUri + "/";
    const template = await readUrlOrPath(host, resolveRelativeUrlOrPath(baseDir, file.path));
    const content = render(template.text, context);
    const destinationFilePath = joinPaths(config.directory, file.destination);
    // create folders in case they don't exist
    await host.mkdirp(getDirectoryPath(destinationFilePath) + "/");
    return host.writeFile(joinPaths(config.directory, file.destination), content);
}
async function getPackageVersion(packageInfo) {
    // TODO: Resolve 'latest' version from npm, issue #1919
    return packageInfo.version ?? "latest";
}
//# sourceMappingURL=scaffold.js.map