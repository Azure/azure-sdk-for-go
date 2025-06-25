import { createDiagnosticCollector } from "../core/diagnostics.js";
import { createDiagnostic } from "../core/messages.js";
import { getDirectoryPath, normalizePath } from "../core/path-utils.js";
import { NoTarget } from "../core/types.js";
import { doIO } from "../utils/io.js";
import { deepClone, omitUndefined } from "../utils/misc.js";
import { expandConfigVariables } from "./config-interpolation.js";
import { loadTypeSpecConfigForPath, validateConfigPathsAbsolute } from "./config-loader.js";
/**
 * Resolve the compiler options for the given entrypoint by resolving the tspconfig.yaml.
 * @param host Compiler host
 * @param options
 */
export async function resolveCompilerOptions(host, options) {
    const diagnostics = createDiagnosticCollector();
    const entrypointStat = await doIO(host.stat, options.entrypoint, (diag) => diagnostics.add(diag), { allowFileNotFound: true });
    const configPath = options.configPath ??
        (entrypointStat?.isDirectory() ? options.entrypoint : getDirectoryPath(options.entrypoint));
    const config = await loadTypeSpecConfigForPath(host, configPath, options.configPath !== undefined, options.configPath === undefined);
    config.diagnostics.forEach((x) => diagnostics.add(x));
    const compilerOptions = diagnostics.pipe(resolveOptionsFromConfig(config, options));
    return diagnostics.wrap(compilerOptions);
}
/**
 * Resolve the compiler options from the given raw TypeSpec config
 * @param config TypeSpec config.
 * @param options Options for interpolation in the config.
 * @returns
 */
export function resolveOptionsFromConfig(config, options) {
    const cwd = normalizePath(options.cwd);
    const diagnostics = createDiagnosticCollector();
    validateConfigNames(config).forEach((x) => diagnostics.add(x));
    const configWithOverrides = {
        ...config,
        ...options.overrides,
        options: mergeOptions(config.options, options.overrides?.options),
    };
    const expandedConfig = diagnostics.pipe(expandConfigVariables(configWithOverrides, {
        cwd,
        outputDir: options.overrides?.outputDir,
        env: options.env ?? {},
        args: options.args,
    }));
    validateConfigPathsAbsolute(expandedConfig).forEach((x) => diagnostics.add(x));
    const resolvedOptions = omitUndefined({
        outputDir: expandedConfig.outputDir,
        config: config.filename,
        configFile: config,
        additionalImports: expandedConfig["imports"],
        warningAsError: expandedConfig.warnAsError,
        trace: expandedConfig.trace,
        emit: expandedConfig.emit,
        options: expandedConfig.options,
        linterRuleSet: expandedConfig.linter,
    });
    return diagnostics.wrap(resolvedOptions);
}
export function validateConfigNames(config) {
    const diagnostics = [];
    function checkName(name) {
        if (name.includes(".")) {
            diagnostics.push(createDiagnostic({
                code: "config-invalid-name",
                format: { name },
                target: NoTarget,
            }));
        }
    }
    function validateNamesRecursively(obj) {
        for (const [key, value] of Object.entries(obj ?? {})) {
            checkName(key);
            if (hasNestedObjects(value)) {
                validateNamesRecursively(value);
            }
        }
    }
    validateNamesRecursively(config.options);
    validateNamesRecursively(config.parameters);
    return diagnostics;
}
function mergeOptions(base, overrides) {
    const configuredEmitters = deepClone(base ?? {});
    function isObject(item) {
        return item && typeof item === "object" && !Array.isArray(item);
    }
    function deepMerge(target, source) {
        if (isObject(target) && isObject(source)) {
            for (const key in source) {
                if (isObject(source[key])) {
                    if (!target[key])
                        Object.assign(target, { [key]: {} });
                    deepMerge(target[key], source[key]);
                }
                else {
                    Object.assign(target, { [key]: source[key] });
                }
            }
        }
        return target;
    }
    for (const [emitterName, cliOptionOverride] of Object.entries(overrides ?? {})) {
        configuredEmitters[emitterName] = deepMerge(configuredEmitters[emitterName] ?? {}, cliOptionOverride);
    }
    return configuredEmitters;
}
function hasNestedObjects(value) {
    return value && typeof value === "object" && !Array.isArray(value);
}
//# sourceMappingURL=config-to-options.js.map