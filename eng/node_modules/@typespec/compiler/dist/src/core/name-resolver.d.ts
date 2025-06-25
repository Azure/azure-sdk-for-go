/**
 * The name resolver is responsible for resolving identifiers to symbols and
 * creating symbols for types that become known during this process. After name
 * resolution, we can do some limited analysis of the reference graph in order
 * to support e.g. augment decorators.
 *
 * Name resolution does not alter any AST nodes or attached symbols in order to
 * ensure AST nodes and attached symbols can be trivially reused between
 * compilations. Instead, symbols created here are either stored in augmented
 * symbol tables or as merged symbols. Any metadata about symbols and nodes are
 * stored in symbol links and node links respectively. The resolver provides
 * APIs for managing this metadata which is useful during later phases.
 *
 * While we resolve some identifiers to symbols during this phase, we often
 * cannot say for sure that an identifier does not exist. Some symbols must be
 * late-bound because the symbol does not become known until after the program
 * has been checked. A common example is members of a model template which often
 * cannot be known until the template is instantiated. Instead, we mark that the
 * reference is unknown and will resolve the symbol (or report an error if it
 * doesn't exist) in later phases. These unknown references cannot be used as
 * the target of an augment decorator.
 *
 * There are some errors we can detect because we have complete symbol
 * information, but we do not report them from here. For example, because we
 * know all namespace bindings and all the declarations inside of them, we could
 * in principle report an error when we attempt to `using` something that isn't
 * a namespace. However, giving a good error message sometimes requires knowing
 * what type was mistakenly referenced, so we merely mark that resolution has
 * failed and move on. Even in cases where we could give a good error we chose
 * not to in order to uniformly handle error reporting in the checker.
 *
 * Name resolution has three sub-phases:
 *
 * 1. Merge namespace symbols and decorator implementation/declaration symbols
 * 2. Resolve using references to namespaces and create namespace-local bindings
 *    for used symbols
 * 3. Resolve type references and bind members
 *
 * The reference resolution and member binding phase implements a deferred
 * resolution strategy. Often we cannot resolve a reference without binding
 * members, but we often cannot bind members without resolving references. In
 * such situations, we stop resolving or binding the current reference or type
 * and attempt to resolve or bind the reference or type it depends on. Once we
 * have done so, we return to the original reference or type and complete our
 * work.
 *
 * This is accomplished by doing a depth-first traversal of the reference graph.
 * On the way down, we discover any dependencies that need to be resolved or
 * bound for the current node, and recurse into the AST nodes, so that on the
 * way back up, all of our dependencies are bound and resolved and we can
 * complete. So while we start with a depth-first traversal of the ASTs in order
 * to discover work to do, most of the actual work is done while following the
 * reference graph, binding and resolving along the way. Circular references are
 * discovered during the reference graph walk and marked as such. Symbol and
 * node links are used to ensure we never resolve the same reference twice. The
 * checker implements a very similar algorithm to evaluate the types of the
 * program.
 **/
import { Mutable } from "../utils/misc.js";
import { Program } from "./program.js";
import { AugmentDecoratorStatementNode, IdentifierNode, MemberExpressionNode, Node, NodeLinks, ResolutionResult, Sym, SymbolLinks, SymbolTable, TypeReferenceNode, UsingStatementNode } from "./types.js";
export interface NameResolver {
    /**
     * Resolve all static symbol links in the program.
     */
    resolveProgram(): void;
    /**
     * Get the merged symbol or itself if not merged.
     * This is the case for Namespace which have multiple nodes and symbol but all reference the same merged one.
     */
    getMergedSymbol(sym: Sym): Sym;
    /**
     * Get augmented symbol table.
     */
    getAugmentedSymbolTable(table: SymbolTable): Mutable<SymbolTable>;
    /**
     * Get node links for the given syntax node.
     * This returns links to which symbol the node reference if applicable(TypeReference, Identifier nodes)
     */
    getNodeLinks(node: Node): NodeLinks;
    /** Get symbol links for the given symbol */
    getSymbolLinks(symbol: Sym): SymbolLinks;
    /** Return augment decorator nodes that are bound to this symbol */
    getAugmentDecoratorsForSym(symbol: Sym): AugmentDecoratorStatementNode[];
    /**
     * Resolve the member expression using the given symbol as base.
     * This can be used to follow the name resolution for template instance which are not statically linked.
     */
    resolveMemberExpressionForSym(sym: Sym, node: MemberExpressionNode, options?: ResolveTypReferenceOptions): ResolutionResult;
    /** Get the meta member by name */
    resolveMetaMemberByName(sym: Sym, name: string): ResolutionResult;
    /** Resolve the given type reference. This should only need to be called on dynamically created nodes that want to resolve which symbol they reference */
    resolveTypeReference(node: TypeReferenceNode | IdentifierNode | MemberExpressionNode): ResolutionResult;
    /** Get the using statement nodes which is not used in resolving yet */
    getUnusedUsings(): UsingStatementNode[];
    /** Built-in symbols. */
    readonly symbols: {
        /** Symbol for the global namespace */
        readonly global: Sym;
        /** Symbol for the null type */
        readonly null: Sym;
    };
}
interface ResolveTypReferenceOptions {
    resolveDecorators?: boolean;
}
export declare function createResolver(program: Program): NameResolver;
export {};
//# sourceMappingURL=name-resolver.d.ts.map