import { setDocData } from "../../core/intrinsic-type-state.js";
const indexTypeKey = Symbol.for(`TypeSpec.index`);
export const indexerDecorator = (context, target, key, value) => {
    const indexer = { key, value };
    context.program.stateMap(indexTypeKey).set(target, indexer);
};
export function getIndexer(program, target) {
    return program.stateMap(indexTypeKey).get(target);
}
/**
 * @internal to be used to set the `@doc` from doc comment.
 */
export const docFromCommentDecorator = (context, target, key, text) => {
    setDocData(context.program, target, key, { value: text, source: "comment" });
};
const prototypeGetterKey = Symbol.for(`TypeSpec.Prototypes.getter`);
/** @internal */
export function getterDecorator(context, target) {
    context.program.stateMap(prototypeGetterKey).set(target, true);
}
/** @internal */
export function isPrototypeGetter(program, target) {
    return program.stateMap(prototypeGetterKey).get(target) ?? false;
}
//# sourceMappingURL=decorators.js.map