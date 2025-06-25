// noop logger shouldn't be used in browser
export function createConsoleSink(options) {
    function log(data) {
        // eslint-disable-next-line no-console
        console.log(formatLog(data));
    }
    return {
        log,
    };
}
export function formatLog(log) {
    return JSON.stringify(log);
}
export function formatDiagnostic(log) {
    return JSON.stringify(log);
}
//# sourceMappingURL=console-sink.browser.js.map