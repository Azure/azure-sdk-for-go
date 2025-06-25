import isUnicodeSupported from "is-unicode-supported";
import pc from "picocolors";
const StatusIcons = {
    success: pc.green("✔"),
    failure: pc.red("×"),
    warn: pc.yellow("⚠"),
    skipped: pc.gray("•"),
};
export class DynamicTask {
    #stream;
    #message;
    #spinner;
    #interval;
    #isTTY;
    #running;
    #finalMessage;
    constructor(message, finalMessage, stream) {
        this.#message = message;
        this.#finalMessage = finalMessage;
        this.#stream = stream;
        this.#spinner = createSpinner();
        this.#isTTY = stream.isTTY && !process.env.CI;
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
        this.#printProgress();
    }
    start() {
        if (this.#isTTY) {
            this.#interval = setInterval(() => {
                this.#printProgress();
            }, 100);
        }
        else {
            this.#stream.write(`- ${this.#message}\n`);
        }
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
        if (this.#interval) {
            clearInterval(this.#interval);
            this.#interval = undefined;
        }
        this.#clear();
        this.#stream.write(`${StatusIcons[status]} ${this.#message}\n`);
    }
    #printProgress() {
        if (!this.#isTTY) {
            return;
        }
        this.#clear();
        this.#clear();
        this.#stream.write(`${pc.yellow(this.#spinner())} ${this.#message}`);
    }
    #clear() {
        if (!this.#isTTY) {
            return;
        }
        this.#stream.cursorTo(0);
        this.#stream.clearLine(0);
    }
}
function createSpinner() {
    let index = 0;
    return () => {
        index = ++index % spinnerFrames.length;
        return spinnerFrames[index];
    };
}
export const spinnerFrames = isUnicodeSupported()
    ? ["⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"]
    : ["-", "\\", "|", "/"];
//# sourceMappingURL=dynamic-task.js.map