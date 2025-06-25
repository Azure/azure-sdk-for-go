import { NpmPackage, NpmPackageProvider } from "./npm-package-provider.js";
export declare class LibraryProvider {
    private npmPackageProvider;
    private filter;
    private libPackageFilterResultCache;
    constructor(npmPackageProvider: NpmPackageProvider, filter: (libExports: Record<string, any>) => boolean);
    /**
     *
     * @param startFolder folder starts to search for package.json with library defined as dependencies
     * @returns
     */
    listLibraries(startFolder: string): Promise<Record<string, NpmPackage>>;
    /**
     *
     * @param startFolder folder starts to search for package.json with library defined as dependencies
     * @param libName
     * @returns
     */
    getLibrary(startFolder: string, libName: string): Promise<NpmPackage | undefined>;
    private getLibFilterResult;
    private getLibraryFromDep;
}
//# sourceMappingURL=lib-provider.d.ts.map