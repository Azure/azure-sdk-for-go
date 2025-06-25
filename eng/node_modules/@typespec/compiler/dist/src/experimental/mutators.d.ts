import { Program } from "../core/program.js";
import { Decorator, FunctionParameter, IntrinsicType, Namespace, TemplateParameter, Type } from "../core/types.js";
import { Realm } from "./realm.js";
/**
 * A description of how a specific kind of type should be mutated.
 *
 * This can either be an object specifying an optional `filter` function and one of `mutate` or `replace`, or simply a
 * function that mutates the type.
 *
 * If a function is provided, it is equivalent to providing an object with a `mutate` function and no `filter` function.
 *
 * @experimental
 */
export type MutatorRecord<T extends Type> = MutatorReplaceRecord<T> | MutatorMutateRecord<T> | MutatorFn<T>;
/**
 * Common functionality for mutator records.
 *
 * @experimental
 */
export interface MutatorRecordCommon<T extends Type> {
    /**
     * A filter function that determines if the mutator should be applied to the type.
     */
    filter?: MutatorFilterFn<T>;
}
/**
 * A mutator that replaces a type's clone with a new type instance.
 *
 * @experimental
 */
export interface MutatorReplaceRecord<T extends Type> extends MutatorRecordCommon<T> {
    /**
     * A mutator function that returns a new type instance to replace the cloned type instance.
     */
    replace: MutatorReplaceFn<T>;
}
/**
 * A mutator that changes the clone of a type in place.
 *
 * @experimental
 */
export interface MutatorMutateRecord<T extends Type> extends MutatorRecordCommon<T> {
    /**
     * A mutator function that edits the clone of the type in place.
     */
    mutate: MutatorFn<T>;
}
/**
 * Edits the clone of the type in place. This function _SHOULD NOT_ modify the source type.
 *
 * @see {@link mutateSubgraph}
 *
 * @param sourceType - The source type.
 * @param clone - The clone of the source type to mutate.
 * @param program - The program in which the `sourceType` occurs.
 * @param realm - The realm in which the `clone` resides.
 *
 * @experimental
 */
export type MutatorFn<T extends Type> = (sourceType: T, clone: T, program: Program, realm: Realm) => void;
/**
 * Determines if the mutator should be applied to the type.
 *
 * This function may either return a boolean or {@link MutatorFlow} flags:
 *
 * - If `true`, the mutator will be applied and will recur (equivalent to `MutatorFlags.MutateAndRecur`).
 * - If `false`, the mutator will not be applied and will recur (equivalent to `MutatorFlags.DoNotMutate`).
 *
 * This predicate runs before the type is cloned.
 *
 * @param sourceType - The source type.
 * @param program - The program in which the `sourceType` occurs.
 * @param realm - The realm where the `sourceType` will be cloned, if this type is mutated.
 *
 * @returns a boolean or {@link MutatorFlow} flags.
 *
 * @experimental
 */
export type MutatorFilterFn<T extends Type> = (sourceType: T, program: Program, realm: Realm) => boolean | MutatorFlow;
/**
 * A function that replaces a mutable type with a new type instance.
 *
 * Returning `clone` from this function is equivalent to providing a `mutate` function instead of a `replace` function.
 *
 * This function runs after the `sourceType` is cloned within the realm.
 *
 * @param sourceType - The source type.
 * @param clone - The clone of the source type to mutate.
 * @param program - The program in which the `sourceType` occurs.
 * @param realm - The realm in which the `clone` resides.
 *
 * @returns a new type instance to replace the cloned type instance.
 *
 * @experimental
 */
export type MutatorReplaceFn<T extends Type> = (sourceType: T, clone: T, program: Program, realm: Realm) => Type;
/**
 * Mutators describe procedures for mutating types in the type graph.
 *
 * Each entry in the mutator describes how to mutate a specific type of node.
 *
 * See {@link mutateSubgraph}.
 *
 * @experimental
 */
export type Mutator = {
    /**
     * The name of this mutator.
     */
    name: string;
} & {
    [Kind in MutableType["kind"]]?: MutatorRecord<Extract<MutableType, {
        kind: Kind;
    }>>;
};
/**
 * A mutator that can additionally mutate namespaces.
 *
 * @experimental
 */
export type MutatorWithNamespace = Mutator & Partial<NamespaceMutator>;
type NamespaceMutator = {
    Namespace?: MutatorRecord<Namespace>;
};
/**
 * Flow control for mutators.
 *
 * When filtering types in a mutator, the filter function may return MutatorFlow flags to control how mutation should
 * proceed.
 *
 * @see {@link MutatorFilterFn}
 *
 * @experimental
 */
export declare enum MutatorFlow {
    /**
     * Mutate the type and recur, further mutating the type's children. This is the default behavior.
     */
    MutateAndRecur = 0,
    /**
     * If this flag is set, the type will not be mutated.
     */
    DoNotMutate = 1,
    /**
     * If this flag is set, the mutator will not proceed recursively into the children of the type.
     */
    DoNotRecur = 2
}
/**
 * A type that can be mutated.
 *
 * @see {@link mutateSubgraph}
 *
 * @experimental
 */
export type MutableType = Exclude<Type, TemplateParameter | IntrinsicType | Decorator | FunctionParameter | Namespace>;
/**
 * Determines if a type is mutable.
 *
 * @experimental
 */
export declare function isMutableType(type: Type): type is MutableType;
/**
 * A mutable type, inclusive of namespaces.
 *
 * @experimental
 */
export type MutableTypeWithNamespace = MutableType | Namespace;
/**
 * Mutate the type graph, allowing namespaces to be mutated.
 *
 * **Warning**: This function will likely mutate the entire type graph. Most TypeSpec types relate to namespaces
 * in some way (e.g. through namespace parent links, or the `namespace` property of a Model).
 *
 * @param program - The program in which the `type` occurs.
 * @param mutators - An array of mutators to apply to the type graph rooted at `type`.
 * @param type - The type to mutate.
 *
 * @returns an object containing the mutated `type` and a nullable `Realm` in which the mutated type resides.
 *
 * @see {@link mutateSubgraph}
 *
 * @experimental
 */
export declare function mutateSubgraphWithNamespace(program: Program, mutators: MutatorWithNamespace[], type: Namespace): {
    realm: Realm | null;
    type: MutableTypeWithNamespace;
};
/**
 * Mutate the type graph.
 *
 * Mutators clone the input `type`, creating a new type instance that is mutated in place.
 *
 * The mutator returns the mutated type and optionally a `realm` in which the mutated clone resides.
 *
 * @see {@link Mutator}
 * @see {@link Realm}
 *
 * **Warning**: Mutators _SHOULD NOT_ modify the source type. Modifications to the source type
 * will be visible to other emitters or libraries that view the original source type, and will
 * be sensitive to the order in which the mutator was applied. Only edit the `clone` type.
 * Furthermore, mutators must take care not to modify elements of the source and clone types
 * that are shared between the two types, such as the properties of any parent references
 * or the `decorators` of the type without taking care to clone them first.
 *
 * @param program - The program in which the `type` occurs.
 * @param mutators - An array of mutators to apply to the type graph rooted at `type`.
 * @param type - The type to mutate.
 *
 * @returns an object containing the mutated `type` and a nullable `Realm` in which the mutated type resides.
 *
 * @experimental
 */
export declare function mutateSubgraph<T extends MutableType>(program: Program, mutators: Mutator[], type: T): {
    realm: Realm | null;
    type: MutableType;
};
export {};
//# sourceMappingURL=mutators.d.ts.map