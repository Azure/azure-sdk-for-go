export class ServerTrackActionTask {
    #log;
    #message;
    #running;
    #finalMessage;
    constructor(message, finalMessage, log) {
        this.#message = message;
        this.#finalMessage = finalMessage;
        this.#log = log;
        this.#running = true;
    }
    get message() {
        return this.#message;
    }
    get isStopped() {
        return !this.#running;
    }
    set message(newMessage) {
        this.#message = newMessage;
    }
    start() {
        this.#log({
            level: "info",
            message: this.#message,
        });
    }
    succeed(message) {
        this.stop("success", message);
    }
    fail(message) {
        this.stop("failure", message);
    }
    warn(message) {
        this.stop("warn", message);
    }
    skip(message) {
        this.stop("skipped", message);
    }
    stop(status, message) {
        this.#running = false;
        this.#message = message ?? this.#finalMessage;
        this.#log({
            level: status !== "failure" ? "info" : "error",
            message: `[${TaskStatusText[status]}] ${this.#message}\n`,
        });
    }
}
const TaskStatusText = {
    success: "succeed",
    failure: "failed",
    warn: "succeeded with warnings",
    skipped: "skipped",
};
/** internal */
export async function trackActionFunc(log, message, finalMessage, asyncAction) {
    const task = new ServerTrackActionTask(message, finalMessage, log);
    task.start();
    try {
        const result = await asyncAction(task);
        if (!task.isStopped) {
            task.succeed();
        }
        return result;
    }
    catch (error) {
        task.fail(message);
        throw error;
    }
}
//# sourceMappingURL=server-track-action-task.js.map