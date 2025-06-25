import { FileEvent } from "vscode-languageserver";
import { CompilerHost } from "../core/types.js";
import { PackageJson } from "../index.js";
export declare class NpmPackageProvider {
    private host;
    private pkgCache;
    constructor(host: CompilerHost);
    notify(changes: FileEvent[]): void;
    /**
     * Search for the nearest package.json file starting from the given folder to its parent/grandparent/... folders
     * @param startFolder the folder to start searching for package.json file
     * @returns
     */
    getPackageJsonFolder(startFolder: string): Promise<string | undefined>;
    /**
     * Get the NpmPackage instance from the folder containing the package.json file.
     *
     * @param packageJsonFolder the dir containing the package.json file. This method won't search for the package.json file, use getPackageJsonFolder to search for the folder containing the package.json file if needed.
     * @returns the NpmPackage instance or undefined if no proper package.json file found
     */
    get(packageJsonFolder: string): Promise<NpmPackage | undefined>;
    private resetCache;
    /**
     * reset the status of the provider with all the caches properly cleaned up
     */
    reset(): void;
}
export declare class NpmPackage {
    private host;
    private packageJsonFolder;
    private packageJsonData;
    private constructor();
    getPackageJsonData(): Promise<PackageJson | undefined>;
    private packageModule;
    getModuleExports(): Promise<Record<string, any> | undefined>;
    resetCache(): void;
    /**
     * Create a NpmPackage instance from a folder containing a package.json file. Make sure to dispose the instance when you finish using it.
     * @param packageJsonFolder the folder containing the package.json file
     * @returns
     */
    static createFrom(host: CompilerHost, packageJsonFolder: string): Promise<NpmPackage | undefined>;
    /**
     *
     * @param packageJsonFolder the folder containing the package.json file
     * @returns
     */
    private static loadNodePackage;
    private static loadModuleExports;
}
//# sourceMappingURL=npm-package-provider.d.ts.map