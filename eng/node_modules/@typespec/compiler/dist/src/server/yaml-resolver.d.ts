import { Position, TextDocument } from "vscode-languageserver-textdocument";
import { Document, Node, Scalar } from "yaml";
import type { TextRange } from "../core/types.js";
import { ServerLog } from "./types.js";
export interface YamlPositionDetail {
    /**
     * The path of the yaml node at the position, it consists of object's property, array's index, empty string for potential object property (empty line)
     */
    path: string[];
    /**
     * The final target of the path, either the key or value of the node pointed by the path property
     *   - "key" or "value" if the node is pointing to an object property
     *   - "arr-item" if the node is pointing to an array item
     */
    type: "key" | "value" | "arr-item";
    /**
     * actual value of target in the doc
     */
    source: string;
    /**
     *  The input quotes (double quotes or single quotes)
     */
    sourceType: Scalar.Type;
    /**
     * The position range of the text in the document, such as [startPos, endPos], see {@link TextRange}
     */
    sourceRange: TextRange | undefined;
    /**
     * The siblings of the target node
     */
    siblings: string[];
    /**
     * The parsed yaml document
     */
    yamlDoc: Document<Node, true>;
    /**
     * The cursor current position
     */
    cursorPosition: number;
}
export declare function resolveYamlPositionDetail(document: TextDocument, position: Position, log: (log: ServerLog) => void): YamlPositionDetail | undefined;
//# sourceMappingURL=yaml-resolver.d.ts.map