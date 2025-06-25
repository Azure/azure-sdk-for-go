const SupportedPackageManagersConfig = {
    npm: {
        commands: {
            install: ["install"],
        },
    },
    // pnpm: {},
    // yarn: {},
};
export function isSupportedPackageManager(value) {
    return value in SupportedPackageManagersConfig;
}
export function getPackageManagerConfig(name) {
    return SupportedPackageManagersConfig[name];
}
//# sourceMappingURL=config.js.map