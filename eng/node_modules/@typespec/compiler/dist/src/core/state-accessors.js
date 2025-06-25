export class StateMap extends Map {
}
export class StateSet extends Map {
}
class StateMapView {
    map;
    constructor(map) {
        this.map = map;
    }
    /** Check if the given type has state */
    has(t) {
        return this.map.has(t) ?? false;
    }
    set(t, v) {
        this.map.set(t, v);
        return this;
    }
    get(t) {
        return this.map.get(t);
    }
    delete(t) {
        return this.map.delete(t);
    }
    /**
     * Danger: Iterating over all types in the state map is not recommended.
     * This occur unexpected result when types are dynamically created, cloned, or removed.
     */
    forEach(cb, thisArg) {
        this.map.forEach(cb, thisArg);
        return this;
    }
    /**
     * Danger: Iterating over all types in the state map is not recommended.
     * This occur unexpected result when types are dynamically created, cloned, or removed.
     */
    get size() {
        return this.map.size;
    }
    clear() {
        return this.map.clear();
    }
    /**
     * Danger: Iterating over all types in the state map is not recommended.
     * This occur unexpected result when types are dynamically created, cloned, or removed.
     */
    entries() {
        return this.map.entries();
    }
    /**
     * Danger: Iterating over all types in the state map is not recommended.
     * This occur unexpected result when types are dynamically created, cloned, or removed.
     */
    values() {
        return this.map.values();
    }
    /**
     * Danger: Iterating over all types in the state map is not recommended.
     * This occur unexpected result when types are dynamically created, cloned, or removed.
     */
    keys() {
        return this.map.keys();
    }
    /**
     * Danger: Iterating over all types in the state map is not recommended.
     * This occur unexpected result when types are dynamically created, cloned, or removed.
     */
    [Symbol.iterator]() {
        return this.entries();
    }
    [Symbol.toStringTag] = "StateMap";
}
class StateSetView {
    set;
    constructor(set) {
        this.set = set;
    }
    has(t) {
        return this.set.has(t) ?? false;
    }
    add(t) {
        this.set.add(t);
        return this;
    }
    delete(t) {
        return this.set.delete(t);
    }
    /**
     * Danger: Iterating over all types in the state map is not recommended.
     * This occur unexpected result when types are dynamically created, cloned, or removed.
     */
    forEach(cb, thisArg) {
        this.set.forEach(cb, thisArg);
        return this;
    }
    get size() {
        return this.set.size;
    }
    clear() {
        return this.set.clear();
    }
    /**
     * Danger: Iterating over all types in the state map is not recommended.
     * This occur unexpected result when types are dynamically created, cloned, or removed.
     */
    values() {
        return this.set.values();
    }
    /**
     * Danger: Iterating over all types in the state map is not recommended.
     * This occur unexpected result when types are dynamically created, cloned, or removed.
     */
    keys() {
        return this.set.keys();
    }
    /**
     * Danger: Iterating over all types in the state map is not recommended.
     * This occur unexpected result when types are dynamically created, cloned, or removed.
     */
    entries() {
        return this.set.entries();
    }
    /**
     * Danger: Iterating over all types in the state map is not recommended.
     * This occur unexpected result when types are dynamically created, cloned, or removed.
     */
    [Symbol.iterator]() {
        return this.values();
    }
    [Symbol.toStringTag] = "StateSet";
}
export function createStateAccessors(stateMaps, stateSets) {
    function stateMap(key) {
        let m = stateMaps.get(key);
        if (!m) {
            m = new Map();
            stateMaps.set(key, m);
        }
        return new StateMapView(m);
    }
    function stateSet(key) {
        let s = stateSets.get(key);
        if (!s) {
            s = new Set();
            stateSets.set(key, s);
        }
        return new StateSetView(s);
    }
    return { stateMap, stateSet };
}
//# sourceMappingURL=state-accessors.js.map