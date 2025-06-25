import { useStateMap } from "../utils/index.js";
import { createStateSymbol } from "./utils.js";
const [getKey, setKey] = useStateMap(createStateSymbol("key"));
export function isKey(program, property) {
    return getKey(program, property) !== undefined;
}
export { getKey as getKeyName, setKey };
//# sourceMappingURL=key.js.map