import type { Type } from "./types.js";
export declare class StateMap extends Map<undefined, Map<Type, unknown>> {
}
export declare class StateSet extends Map<undefined, Set<Type>> {
}
declare class StateMapView<V> implements Map<Type, V> {
    private map;
    constructor(map: Map<Type, V>);
    /** Check if the given type has state */
    has(t: Type): boolean;
    set(t: Type, v: any): this;
    get(t: Type): V | undefined;
    delete(t: Type): boolean;
    /**
     * Danger: Iterating over all types in the state map is not recommended.
     * This occur unexpected result when types are dynamically created, cloned, or removed.
     */
    forEach(cb: (value: V, key: Type, map: Map<Type, V>) => void, thisArg?: any): this;
    /**
     * Danger: Iterating over all types in the state map is not recommended.
     * This occur unexpected result when types are dynamically created, cloned, or removed.
     */
    get size(): number;
    clear(): void;
    /**
     * Danger: Iterating over all types in the state map is not recommended.
     * This occur unexpected result when types are dynamically created, cloned, or removed.
     */
    entries(): MapIterator<[Type, V]>;
    /**
     * Danger: Iterating over all types in the state map is not recommended.
     * This occur unexpected result when types are dynamically created, cloned, or removed.
     */
    values(): MapIterator<V>;
    /**
     * Danger: Iterating over all types in the state map is not recommended.
     * This occur unexpected result when types are dynamically created, cloned, or removed.
     */
    keys(): MapIterator<Type>;
    /**
     * Danger: Iterating over all types in the state map is not recommended.
     * This occur unexpected result when types are dynamically created, cloned, or removed.
     */
    [Symbol.iterator](): MapIterator<[Type, V]>;
    [Symbol.toStringTag]: string;
}
declare class StateSetView implements Set<Type> {
    private set;
    constructor(set: Set<Type>);
    has(t: Type): boolean;
    add(t: Type): this;
    delete(t: Type): boolean;
    /**
     * Danger: Iterating over all types in the state map is not recommended.
     * This occur unexpected result when types are dynamically created, cloned, or removed.
     */
    forEach(cb: (value: Type, value2: Type, set: Set<Type>) => void, thisArg?: any): this;
    get size(): number;
    clear(): void;
    /**
     * Danger: Iterating over all types in the state map is not recommended.
     * This occur unexpected result when types are dynamically created, cloned, or removed.
     */
    values(): SetIterator<Type>;
    /**
     * Danger: Iterating over all types in the state map is not recommended.
     * This occur unexpected result when types are dynamically created, cloned, or removed.
     */
    keys(): SetIterator<Type>;
    /**
     * Danger: Iterating over all types in the state map is not recommended.
     * This occur unexpected result when types are dynamically created, cloned, or removed.
     */
    entries(): SetIterator<[Type, Type]>;
    /**
     * Danger: Iterating over all types in the state map is not recommended.
     * This occur unexpected result when types are dynamically created, cloned, or removed.
     */
    [Symbol.iterator](): SetIterator<Type>;
    [Symbol.toStringTag]: string;
}
export declare function createStateAccessors(stateMaps: Map<symbol, Map<Type, unknown>>, stateSets: Map<symbol, Set<Type>>): {
    stateMap: <T>(key: symbol) => StateMapView<T>;
    stateSet: (key: symbol) => StateSetView;
};
export {};
//# sourceMappingURL=state-accessors.d.ts.map