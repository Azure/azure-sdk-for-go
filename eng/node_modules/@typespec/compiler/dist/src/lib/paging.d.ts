import type { ContinuationTokenDecorator, FirstLinkDecorator, LastLinkDecorator, ListDecorator, NextLinkDecorator, OffsetDecorator, PageIndexDecorator, PageItemsDecorator, PageSizeDecorator, PrevLinkDecorator } from "../../generated-defs/TypeSpec.js";
import { Program } from "../core/program.js";
import type { Diagnostic, ModelProperty, Operation } from "../core/types.js";
export declare const 
/**
 * Check if the given operation is used to page through a list.
 * @param program Program
 * @param target Operation
 */
isList: (program: Program, type: Operation) => boolean, markList: (program: Program, type: Operation) => void, 
/** {@inheritdoc ListDecorator} */
listDecorator: ListDecorator;
export declare const 
/**
 * Check if the given property is the `@offset` property for a paging operation.
 * @param program Program
 * @param target Model Property
 */
isOffsetProperty: (program: Program, type: ModelProperty) => boolean, markOffset: (program: Program, type: ModelProperty) => void, 
/** {@inheritdoc OffsetDecorator} */
offsetDecorator: OffsetDecorator;
export declare const 
/**
 * Check if the given property is the `@pageIndex` property for a paging operation.
 * @param program Program
 * @param target Model Property
 */
isPageIndexProperty: (program: Program, type: ModelProperty) => boolean, markPageIndexProperty: (program: Program, type: ModelProperty) => void, 
/** {@inheritdoc PageIndexDecorator} */
pageIndexDecorator: PageIndexDecorator;
export declare const 
/**
 * Check if the given property is the `@pageIndex` property for a paging operation.
 * @param program Program
 * @param target Model Property
 */
isPageSizeProperty: (program: Program, type: ModelProperty) => boolean, markPageSizeProperty: (program: Program, type: ModelProperty) => void, 
/** {@inheritdoc PageSizeDecorator} */
pageSizeDecorator: PageSizeDecorator;
export declare const 
/**
 * Check if the given property is the `@pageIndex` property for a paging operation.
 * @param program Program
 * @param target Model Property
 */
isPageItemsProperty: (program: Program, type: ModelProperty) => boolean, markPageItemsProperty: (program: Program, type: ModelProperty) => void, 
/** {@inheritdoc PageItemsDecorator} */
pageItemsDecorator: PageItemsDecorator;
export declare const 
/**
 * Check if the given property is the `@pageIndex` property for a paging operation.
 * @param program Program
 * @param target Model Property
 */
isContinuationTokenProperty: (program: Program, type: ModelProperty) => boolean, markContinuationTokenProperty: (program: Program, type: ModelProperty) => void, 
/** {@inheritdoc ContinuationTokenDecorator} */
continuationTokenDecorator: ContinuationTokenDecorator;
export declare const 
/**
 * Check if the given property is the `@nextLink` property for a paging operation.
 * @param program Program
 * @param target Model Property
 */
isNextLink: (program: Program, type: ModelProperty) => boolean, markNextLink: (program: Program, type: ModelProperty) => void, 
/** {@inheritdoc NextLinkDecorator} */
nextLinkDecorator: NextLinkDecorator;
export declare const 
/**
 * Check if the given property is the `@prevLink` property for a paging operation.
 * @param program Program
 * @param target Model Property
 */ isPrevLink: (program: Program, type: ModelProperty) => boolean, markPrevLink: (program: Program, type: ModelProperty) => void, 
/** {@inheritdoc PrevLinkDecorator} */
prevLinkDecorator: PrevLinkDecorator;
export declare const 
/**
 * Check if the given property is the `@firstLink` property for a paging operation.
 * @param program Program
 * @param target Model Property
 */ isFirstLink: (program: Program, type: ModelProperty) => boolean, markFirstLink: (program: Program, type: ModelProperty) => void, 
/** {@inheritdoc FirstLinkDecorator} */
firstLinkDecorator: FirstLinkDecorator;
export declare const 
/**
 * Check if the given property is the `@lastLink` property for a paging operation.
 * @param program Program
 * @param target Model Property
 */ isLastLink: (program: Program, type: ModelProperty) => boolean, markLastLink: (program: Program, type: ModelProperty) => void, 
/** {@inheritdoc LastLinkDecorator} */
lastLinkDecorator: LastLinkDecorator;
export declare function validatePagingOperations(program: Program): void;
export interface PagingProperty {
    readonly property: ModelProperty;
    /**
     * If the paging property is nested, this will contain the path to the paging property in the model
     * and array length will be greater than one.
     *
     * You can use this to generate the path to the property in the model with following code:
     * @example
     * ```ts
     * const path = pagingInfo.output.pageItems.path.map((prop) => prop.name).join(".");
     * ```
     */
    readonly path: ModelProperty[];
}
export interface PagingOperation {
    readonly input: {
        readonly offset?: PagingProperty;
        readonly pageIndex?: PagingProperty;
        readonly pageSize?: PagingProperty;
        readonly continuationToken?: PagingProperty;
    };
    readonly output: {
        readonly pageItems: PagingProperty;
        readonly nextLink?: PagingProperty;
        readonly prevLink?: PagingProperty;
        readonly firstLink?: PagingProperty;
        readonly lastLink?: PagingProperty;
        readonly continuationToken?: PagingProperty;
    };
}
export declare function getPagingOperation(program: Program, op: Operation): [PagingOperation | undefined, readonly Diagnostic[]];
//# sourceMappingURL=paging.d.ts.map