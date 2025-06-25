import { TaskStatus, TrackActionTask } from "../types.js";
export declare class DynamicTask implements TrackActionTask {
    #private;
    constructor(message: string, finalMessage: string, stream: NodeJS.WriteStream);
    get message(): string;
    get isStopped(): boolean;
    set message(newMessage: string);
    start(): void;
    succeed(message?: string): void;
    fail(message?: string): void;
    warn(message?: string): void;
    skip(message?: string): void;
    stop(status: TaskStatus, message?: string): void;
}
export declare const spinnerFrames: string[];
//# sourceMappingURL=dynamic-task.d.ts.map