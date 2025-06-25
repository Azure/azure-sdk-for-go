import { TaskStatus, TrackActionTask } from "../core/types.js";
import { ServerLog } from "./types.js";
export declare class ServerTrackActionTask implements TrackActionTask {
    #private;
    constructor(message: string, finalMessage: string, log: (log: ServerLog) => void);
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
/** internal */
export declare function trackActionFunc<T>(log: (log: ServerLog) => void, message: string, finalMessage: string, asyncAction: (task: TrackActionTask) => Promise<T>): Promise<T>;
//# sourceMappingURL=server-track-action-task.d.ts.map