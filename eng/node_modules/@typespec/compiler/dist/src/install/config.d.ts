export interface PackageManagerConfig {
    commands: {
        install: string[];
    };
}
export type SupportedPackageManager = "npm";
export declare function isSupportedPackageManager(value: string): value is SupportedPackageManager;
export declare function getPackageManagerConfig(name: SupportedPackageManager): PackageManagerConfig;
//# sourceMappingURL=config.d.ts.map