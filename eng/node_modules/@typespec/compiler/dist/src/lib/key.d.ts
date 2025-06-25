import { Program } from "../core/program.js";
import { ModelProperty, Type } from "../core/types.js";
declare const getKey: (program: Program, type: Type) => string | undefined, setKey: (program: Program, type: Type, value: string) => void;
export declare function isKey(program: Program, property: ModelProperty): boolean;
export { getKey as getKeyName, setKey };
//# sourceMappingURL=key.d.ts.map