import type { ServiceDecorator } from "../../generated-defs/TypeSpec.js";
import type { Program } from "../core/program.js";
import { Namespace } from "../core/types.js";
export interface ServiceDetails {
    title?: string;
}
export interface Service extends ServiceDetails {
    type: Namespace;
}
declare const getService: (program: Program, type: Namespace) => Service | undefined;
/**
 * List all the services defined in the TypeSpec program
 * @param program Program
 * @returns List of service.
 */
export declare function listServices(program: Program): Service[];
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
export declare function isService(program: Program, namespace: Namespace): boolean;
/**
 * Mark the given namespace as a service.
 * @param program Program
 * @param namespace Namespace
 * @param details Service details
 */
export declare function addService(program: Program, namespace: Namespace, details?: ServiceDetails): void;
export declare const $service: ServiceDecorator;
//# sourceMappingURL=service.d.ts.map