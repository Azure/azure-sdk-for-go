var _a;
import { compilerAssert } from "../core/diagnostics.js";
import { createTypekit } from "./typekit/index.js";
/**
 * A Realm's view of a Program's state map for a given state key.
 *
 * For all operations, if a type was created within the realm, the realm's own state map is used. Otherwise, the owning'
 * Program's state map is used.
 *
 * @experimental
 */
class StateMapRealmView {
    #realm;
    #parentState;
    #realmState;
    constructor(realm, realmState, parentState) {
        this.#realm = realm;
        this.#parentState = parentState;
        this.#realmState = realmState;
    }
    has(t) {
        return this.#select(t).has(t) ?? false;
    }
    set(t, v) {
        this.#select(t).set(t, v);
        return this;
    }
    get(t) {
        return this.#select(t).get(t);
    }
    delete(t) {
        return this.#select(t).delete(t);
    }
    forEach(cb, thisArg) {
        for (const item of this.entries()) {
            cb.call(thisArg, item[1], item[0], this);
        }
        return this;
    }
    get size() {
        return this.#realmState.size + this.#parentState.size;
    }
    clear() {
        this.#realmState.clear();
    }
    *entries() {
        for (const item of this.#realmState) {
            yield item;
        }
        for (const item of this.#parentState) {
            yield item;
        }
        return undefined;
    }
    *values() {
        for (const item of this.entries()) {
            yield item[1];
        }
        return undefined;
    }
    *keys() {
        for (const item of this.entries()) {
            yield item[0];
        }
        return undefined;
    }
    [Symbol.iterator]() {
        return this.entries();
    }
    [Symbol.toStringTag] = "StateMap";
    #select(keyType) {
        if (this.#realm.hasType(keyType)) {
            return this.#realmState;
        }
        return this.#parentState;
    }
}
/**
 * A Realm is an alternate view of a Program where types can be cloned, deleted, and modified without affecting the
 * original types in the Program.
 *
 * The realm stores the types that exist within the realm, views of state maps that only apply within the realm,
 * and a view of types that have been removed from the realm's view.
 *
 * @experimental
 */
export class Realm {
    #program;
    /**
     * Stores all types owned by this realm.
     */
    #types = new Set();
    /**
     * Stores types that are deleted in this realm. When a realm is active and doing a traversal, you will
     * not find this type in e.g. collections. Deleted types are mapped to `null` if you ask for it.
     */
    #deletedTypes = new WeakSet();
    #stateMaps = new Map();
    key;
    /**
     * Create a new realm in the given program.
     *
     * @param program - The program to create the realm in.
     * @param description - A short description of the realm's purpose.
     */
    constructor(program, description) {
        this.key = Symbol(description);
        this.#program = program;
    }
    #_typekit;
    /**
     * The typekit instance bound to this realm.
     *
     * If the realm does not already have a typekit associated with it, one will be created and bound to this realm.
     */
    get typekit() {
        return (this.#_typekit ??= createTypekit(this));
    }
    /**
     * The program that this realm is associated with.
     */
    get program() {
        return this.#program;
    }
    /**
     * Gets a state map for the given state key symbol.
     *
     * This state map is a view of the program's state map for the given state key, with modifications made to the realm's
     * own state.
     *
     * @param stateKey - The symbol to use as the state key.
     * @returns The realm's state map for the given state key.
     */
    stateMap(stateKey) {
        let m = this.#stateMaps.get(stateKey);
        if (!m) {
            m = new Map();
            this.#stateMaps.set(stateKey, m);
        }
        return new StateMapRealmView(this, m, this.#program.stateMap(stateKey));
    }
    /**
     * Clones a type and adds it to the realm. This operation will use the realm's typekit to clone the type.
     *
     * @param type - The type to clone.
     * @returns A clone of the input type that exists within this realm.
     */
    clone(type) {
        compilerAssert(type, "Undefined type passed to clone");
        const clone = this.#cloneIntoRealm(type);
        this.typekit.type.finishType(clone);
        return clone;
    }
    /**
     * Removes a type from this realm. This operation will not affect the type in the program, only this realm's view
     * of the type.
     *
     * @param type - The TypeSpec type to remove from this realm.
     */
    remove(type) {
        this.#deletedTypes.add(type);
    }
    /**
     * Determines whether or not this realm contains a given type.
     *
     * @param type - The type to check.
     * @returns true if the type was created within this realm or added to this realm, false otherwise.
     */
    hasType(type) {
        return this.#types.has(type);
    }
    /**
     * Adds a type to this realm. Once a type is added to the realm, the realm considers it part of itself.
     *
     * A type can be present in multiple realms, but `Realm.realmForType` will only return the last realm that the type
     * was added to.
     *
     * @param type - The type to add to this realm.
     */
    addType(type) {
        this.#types.add(type);
        _a.realmForType.set(type, this);
    }
    #cloneIntoRealm(type) {
        const clone = this.typekit.type.clone(type);
        this.#types.add(clone);
        _a.realmForType.set(clone, this);
        return clone;
    }
    // TODO better way?
    /** @internal */
    get types() {
        return this.#types;
    }
    static realmForType = singleton("Realm.realmForType", () => new WeakMap());
}
_a = Realm;
/**
 * Create a singleton instance that is shared across the process.
 * This is to have a true singleton even if multiple instance of the compiler/library are loaded.
 * @param key - The key to use for the singleton.
 * @param init - The function to call to create the singleton.
 */
function singleton(key, init) {
    const sym = Symbol.for(key);
    if (!globalThis[sym]) {
        globalThis[sym] = init();
    }
    return globalThis[sym];
}
//# sourceMappingURL=realm.js.map