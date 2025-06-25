import { validateDecoratorUniqueOnNode } from "../core/decorator-utils.js";
import { Realm } from "../experimental/realm.js";
import { useStateMap } from "../utils/index.js";
const [getService, setService, getServiceMap] = useStateMap(Symbol.for("@typespec/compiler.services"));
/**
 * List all the services defined in the TypeSpec program
 * @param program Program
 * @returns List of service.
 */
export function listServices(program) {
    return [...getServiceMap(program).values()].filter((x) => !Realm.realmForType.has(x.type));
}
export { 
/**
 * Get the service information for the given namespace.
 * @param program Program
 * @param namespace Service namespace
 * @returns Service information or undefined if namespace is not a service namespace.
 */
getService, };
/**
 * Check if the namespace is defined as a service.
 * @param program Program
 * @param namespace Namespace
 * @returns Boolean
 */
export function isService(program, namespace) {
    return getServiceMap(program).has(namespace);
}
/**
 * Mark the given namespace as a service.
 * @param program Program
 * @param namespace Namespace
 * @param details Service details
 */
export function addService(program, namespace, details = {}) {
    const serviceMap = getServiceMap(program);
    const existing = serviceMap.get(namespace) ?? {};
    setService(program, namespace, { ...existing, ...details, type: namespace });
}
export const $service = (context, target, options) => {
    validateDecoratorUniqueOnNode(context, target, $service);
    addService(context.program, target, options);
};
//# sourceMappingURL=service.js.map