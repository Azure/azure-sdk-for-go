import { realpath } from "fs";
import { mkdir, readdir, readFile, rm, stat, writeFile } from "fs/promises";
import { createSourceFile } from "./source-file.js";
/**
 * Implementation of the @see SystemHost using the real file system.
 */
export const NodeSystemHost = {
    readUrl: async (url) => {
        const response = await fetch(url, { redirect: "follow" });
        const text = await response.text();
        return createSourceFile(text, response.url);
    },
    readFile: async (path) => createSourceFile(await readUtf8File(path), path),
    writeFile: (path, content) => writeFile(path, content, { encoding: "utf-8" }),
    readDir: (path) => readdir(path),
    rm: (path, options) => rm(path, options),
    stat(path) {
        return stat(path);
    },
    realpath(path) {
        // BUG in the promise version of realpath https://github.com/microsoft/typespec/issues/2783
        // Fix was only made to node 21.6 at this time. https://github.com/nodejs/node/issues/51031
        return new Promise((resolve, reject) => {
            realpath(path, (err, resolvedPath) => {
                if (err) {
                    reject(err);
                }
                else {
                    resolve(resolvedPath);
                }
            });
        });
    },
    mkdirp: (path) => mkdir(path, { recursive: true }),
};
async function readUtf8File(path) {
    const buffer = await readFile(path);
    const len = buffer.length;
    if (len >= 2 && buffer[0] === 0xfe && buffer[1] === 0xff) {
        throw new InvalidEncodingError("UTF-16 BE");
    }
    if (len >= 2 && buffer[0] === 0xff && buffer[1] === 0xfe) {
        throw new InvalidEncodingError("UTF-16 LE");
    }
    if (len >= 3 && buffer[0] === 0xef && buffer[1] === 0xbb && buffer[2] === 0xbf) {
        // UTF-8 byte order mark detected
        return buffer.toString("utf8", 3);
    }
    // Default is UTF-8 with no byte order mark
    return buffer.toString("utf8");
}
export class InvalidEncodingError extends Error {
    constructor(encoding) {
        super(`Invalid encoding ${encoding}. TypeSpec only supports UTF-8 and UTF-8 with bom`);
    }
}
//# sourceMappingURL=node-system-host.js.map