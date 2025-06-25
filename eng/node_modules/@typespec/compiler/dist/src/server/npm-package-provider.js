import { getDirectoryPath, joinPaths, normalizePath } from "../core/path-utils.js";
import { loadJsFile } from "../core/source-loader.js";
import { NoTarget } from "../core/types.js";
import { resolveModule } from "../index.js";
import { distinctArray, isWhitespaceStringOrUndefined, tryParseJson } from "../utils/misc.js";
export class NpmPackageProvider {
    host;
    pkgCache = new Map();
    constructor(host) {
        this.host = host;
    }
    notify(changes) {
        let folders = changes
            .map((c) => normalizePath(this.host.fileURLToPath(c.uri)))
            .filter((c) => c.endsWith("/package.json"))
            .map((c) => getDirectoryPath(c));
        folders = distinctArray(folders, (f) => f);
        for (const folder of folders) {
            const pkg = this.pkgCache.get(folder);
            if (pkg) {
                pkg.resetCache();
                // since we may not get the notification for changes under node_modules
                // just reset those for safety
                const nodeModulesFolder = joinPaths(folder, "node_modules");
                this.pkgCache.forEach((nmPkg, key) => {
                    if (key.startsWith(nodeModulesFolder)) {
                        nmPkg.resetCache();
                    }
                });
            }
        }
    }
    /**
     * Search for the nearest package.json file starting from the given folder to its parent/grandparent/... folders
     * @param startFolder the folder to start searching for package.json file
     * @returns
     */
    async getPackageJsonFolder(startFolder) {
        if (isWhitespaceStringOrUndefined(startFolder)) {
            return undefined;
        }
        let lastFolder = "";
        let curFolder = startFolder;
        while (curFolder !== lastFolder) {
            const packageJsonPath = joinPaths(curFolder, "package.json");
            try {
                const stat = await this.host.stat(packageJsonPath);
                if (stat.isFile()) {
                    return curFolder;
                }
            }
            catch (e) {
                // ignore
            }
            lastFolder = curFolder;
            curFolder = getDirectoryPath(curFolder);
        }
        return undefined;
    }
    /**
     * Get the NpmPackage instance from the folder containing the package.json file.
     *
     * @param packageJsonFolder the dir containing the package.json file. This method won't search for the package.json file, use getPackageJsonFolder to search for the folder containing the package.json file if needed.
     * @returns the NpmPackage instance or undefined if no proper package.json file found
     */
    async get(packageJsonFolder) {
        const key = normalizePath(packageJsonFolder);
        const r = this.pkgCache.get(key);
        if (r) {
            return r;
        }
        else {
            const pkg = await NpmPackage.createFrom(this.host, packageJsonFolder);
            if (pkg) {
                this.pkgCache.set(key, pkg);
                return pkg;
            }
            else {
                return undefined;
            }
        }
    }
    resetCache() {
        const t = this.pkgCache;
        this.pkgCache = new Map();
        t.forEach((pkg) => {
            pkg.resetCache();
        });
    }
    /**
     * reset the status of the provider with all the caches properly cleaned up
     */
    reset() {
        this.resetCache();
    }
}
export class NpmPackage {
    host;
    packageJsonFolder;
    packageJsonData;
    constructor(host, packageJsonFolder, packageJsonData) {
        this.host = host;
        this.packageJsonFolder = packageJsonFolder;
        this.packageJsonData = packageJsonData;
    }
    async getPackageJsonData() {
        if (!this.packageJsonData) {
            this.packageJsonData = await NpmPackage.loadNodePackage(this.host, this.packageJsonFolder);
        }
        return this.packageJsonData;
    }
    packageModule;
    async getModuleExports() {
        if (!this.packageModule) {
            const data = await this.getPackageJsonData();
            if (!data)
                return undefined;
            this.packageModule = await NpmPackage.loadModuleExports(this.host, this.packageJsonFolder, data.name);
        }
        return this.packageModule;
    }
    resetCache() {
        this.packageJsonData = undefined;
        this.packageModule = undefined;
    }
    /**
     * Create a NpmPackage instance from a folder containing a package.json file. Make sure to dispose the instance when you finish using it.
     * @param packageJsonFolder the folder containing the package.json file
     * @returns
     */
    static async createFrom(host, packageJsonFolder) {
        if (!packageJsonFolder) {
            return undefined;
        }
        const data = await NpmPackage.loadNodePackage(host, packageJsonFolder);
        if (!data) {
            return undefined;
        }
        return new NpmPackage(host, packageJsonFolder, data);
    }
    /**
     *
     * @param packageJsonFolder the folder containing the package.json file
     * @returns
     */
    static async loadNodePackage(host, packageJsonFolder) {
        if (!packageJsonFolder) {
            return undefined;
        }
        const packageJsonPath = joinPaths(packageJsonFolder, "package.json");
        try {
            if (!(await host.stat(packageJsonPath)).isFile()) {
                return undefined;
            }
            const content = await host.readFile(packageJsonPath);
            const data = tryParseJson(content.text);
            if (!data || !data.name) {
                return undefined;
            }
            return data;
        }
        catch {
            return undefined;
        }
    }
    static async loadModuleExports(host, baseDir, packageName) {
        try {
            const module = await resolveModule({
                realpath: host.realpath,
                readFile: async (path) => {
                    const sf = await host.readFile(path);
                    return sf.text;
                },
                stat: host.stat,
            }, packageName, { baseDir });
            if (!module) {
                return undefined;
            }
            const entrypoint = module.type === "file" ? module.path : module.mainFile;
            const oldExit = process.exit;
            try {
                // override process.exit to prevent the process from exiting because of it's called in loaded js file
                let result;
                process.exit = (() => {
                    // for module that calls process.exit when being imported, create an empty object as it's exports to avoid load it again
                    result = {};
                    throw new Error("process.exit is called unexpectedly when loading js file: " + entrypoint);
                });
                const [file] = await loadJsFile(host, entrypoint, NoTarget);
                return result ?? file?.esmExports;
            }
            finally {
                process.exit = oldExit;
            }
        }
        catch (e) {
            return undefined;
        }
    }
}
//# sourceMappingURL=npm-package-provider.js.map