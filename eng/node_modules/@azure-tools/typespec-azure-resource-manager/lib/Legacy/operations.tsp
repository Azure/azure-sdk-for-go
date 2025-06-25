import "@typespec/rest";
import "@typespec/http";

namespace Azure.ResourceManager.Legacy;

using Http;
using Azure.ResourceManager.Foundations;
using Azure.ResourceManager.Private;
using Rest;

/**
 * An operation template used to build resource operations in which the same resource type
 * is accessible at multiple, fixed resource paths.
 * @template ParentParameters The path parameters for the resource parent
 * @template ResourceTypeParameter The path parameter for the resource name
 * @template ErrorType Optional. The type of error models used in operations created form this template
 */
@doc("")
interface LegacyOperations<
  ParentParameters extends {},
  ResourceTypeParameter extends {},
  ErrorType extends {} = ErrorResponse
> {
  /**
   * A long-running resource CreateOrUpdate (PUT)
   * @template Resource the resource being created or updated
   * @template LroHeaders Optional.  Allows overriding the lro headers returned on resource create
   * @template Parameters Optional. Additional parameters after the path parameters
   * @template Response Optional. The success response(s) for the PUT operation
   * @template OptionalRequestBody Optional. Indicates whether the request body is optional
   */
  @autoRoute
  @doc("Create a {name}", Resource)
  @armResourceCreateOrUpdate(Resource)
  @Azure.Core.Foundations.Private.defaultFinalStateVia(#["location", "azure-async-operation"])
  @put
  CreateOrUpdateAsync<
    Resource extends Foundations.SimpleResource,
    LroHeaders extends TypeSpec.Reflection.Model = ArmAsyncOperationHeader<FinalResult = Resource> &
      Azure.Core.Foundations.RetryAfterHeader,
    Parameters extends {} = {},
    Response extends {} = ArmResourceUpdatedResponse<Resource> | ArmResourceCreatedResponse<
      Resource,
      LroHeaders
    >,
    OptionalRequestBody extends valueof boolean = false
  >(
    ...ParentParameters,
    ...ResourceTypeParameter,
    ...Parameters,
    @doc("Resource create parameters.") @armBodyRoot(OptionalRequestBody) resource: Resource,
  ): Response | ErrorType;

  /**
   * A synchronous resource CreateOrUpdate (PUT)
   * @template Resource the resource being created or updated
   * @template Parameters Optional. Additional parameters after the path parameters
   * @template Response Optional. The success response(s) for the PUT operation
   * @template OptionalRequestBody Optional. Indicates whether the request body is optional
   */
  #suppress "@azure-tools/typespec-azure-core/no-private-usage"
  @autoRoute
  @doc("Create a {name}", Resource)
  @armResourceCreateOrUpdate(Resource)
  @put
  CreateOrUpdateSync<
    Resource extends Foundations.SimpleResource,
    Parameters extends {} = {},
    Response extends {} = ArmResourceUpdatedResponse<Resource> | ArmResourceCreatedSyncResponse<Resource>,
    OptionalRequestBody extends valueof boolean = false
  >(
    ...ParentParameters,
    ...ResourceTypeParameter,
    ...Parameters,
    @doc("Resource create parameters.") @armBodyRoot(OptionalRequestBody) resource: Resource,
  ): Response | ErrorType;

  /**
   * A long-running resource Update (PATCH)
   * @template Resource the resource being created or updated
   * @template PatchModel the PATCH request model
   * @template LroHeaders Optional.  Allows overriding the lro headers returned on resource create
   * @template Parameters Optional. Additional parameters after the path parameters
   * @template Response Optional. The success response(s) for the PATCH operation
   * @template OptionalRequestBody Optional. Indicates whether the request body is optional
   */
  @autoRoute
  @doc("Update a {name}", Resource)
  @armResourceUpdate(Resource)
  @patch(#{ implicitOptionality: false })
  CustomPatchAsync<
    Resource extends Foundations.SimpleResource,
    PatchModel extends {} = Azure.ResourceManager.Foundations.TagsUpdateModel<Resource>,
    LroHeaders extends TypeSpec.Reflection.Model = ArmLroLocationHeader<
      Azure.Core.StatusMonitorPollingOptions<ArmOperationStatus>,
      Resource,
      string
    > &
      Azure.Core.Foundations.RetryAfterHeader,
    Parameters extends {} = {},
    Response extends {} = ArmResponse<Resource> | ArmAcceptedLroResponse<
      "Resource update request accepted.",
      LroHeaders
    >,
    OptionalRequestBody extends valueof boolean = false
  >(
    ...ParentParameters,
    ...ResourceTypeParameter,
    ...Parameters,
    @doc("Resource create parameters.") @armBodyRoot(OptionalRequestBody) properties: PatchModel,
  ): Response | ErrorType;

  /**
   * A synchronous resource Update (PATCH)
   * @template Resource the resource being created or updated
   * @template PatchModel the PATCH request model
   * @template Parameters Optional. Additional parameters after the path parameters
   * @template Response Optional. The success response(s) for the PATCH operation
   * @template OptionalRequestBody Optional. Indicates whether the request body is optional
   */
  @autoRoute
  @doc("Update a {name}", Resource)
  @armResourceUpdate(Resource)
  @patch(#{ implicitOptionality: false })
  CustomPatchSync<
    Resource extends Foundations.SimpleResource,
    PatchModel extends {} = Azure.ResourceManager.Foundations.TagsUpdateModel<Resource>,
    Parameters extends {} = {},
    Response extends {} = ArmResponse<Resource>,
    OptionalRequestBody extends valueof boolean = false
  >(
    ...ParentParameters,
    ...ResourceTypeParameter,
    ...Parameters,
    @doc("Resource create parameters.") @armBodyRoot(OptionalRequestBody) properties: PatchModel,
  ): Response | ErrorType;

  /**
   * Delete a resource asynchronously
   * @template Resource The resource being deleted
   * @template LroHeaders The lro headers for the operation
   * @template Parameters Optional. Additional parameters after the path parameters
   * @template Response Optional. The success response(s) for the delete operation
   */
  #suppress "@azure-tools/typespec-azure-core/no-response-body" "Valid"
  @autoRoute
  @doc("Delete a {name}", Resource)
  @delete
  @deletesResource(Resource)
  @armResourceDelete(Resource)
  DeleteWithoutOkAsync<
    Resource extends Foundations.SimpleResource,
    LroHeaders extends TypeSpec.Reflection.Model = ArmLroLocationHeader &
      Azure.Core.Foundations.RetryAfterHeader,
    Parameters extends {} = {},
    Response extends {} = ArmDeleteAcceptedLroResponse<LroHeaders> | ArmDeletedNoContentResponse
  >(...ParentParameters, ...ResourceTypeParameter, ...Parameters): Response | ErrorType;

  /**
   * Delete a resource synchronously
   * @template Resource The resource being deleted
   * @template Parameters Optional. Additional parameters after the path parameters
   * @template Response Optional. The success response(s) for the delete operation
   */
  #suppress "@azure-tools/typespec-azure-core/no-response-body" "Valid"
  @autoRoute
  @doc("Delete a {name}", Resource)
  @delete
  @deletesResource(Resource)
  @armResourceDelete(Resource)
  DeleteSync<
    Resource extends Foundations.SimpleResource,
    Parameters extends {} = {},
    Response extends {} = ArmDeletedResponse | ArmDeletedNoContentResponse
  >(...ParentParameters, ...ResourceTypeParameter, ...Parameters): Response | ErrorType;

  /**
   * @dev Get a resource
   * @template Resource The resource being read
   * @template Parameters Optional. Additional parameters after the path parameters
   * @template Response Optional. The success response for a get operation.
   */
  @autoRoute
  @doc("Get a {name}", Resource)
  @get
  @readsResource(Resource)
  @armResourceRead(Resource)
  Read<
    Resource extends Foundations.SimpleResource,
    Parameters extends {} = {},
    Response extends {} = ArmResponse<Resource>
  >(...ParentParameters, ...ResourceTypeParameter, ...Parameters): Response | ErrorType;

  /**
   * List a resource
   * @template Resource The resource being listed
   * @template Parameters Optional. Additional parameters after the path parameters
   * @template Response Optional. The response returned by the list
   */
  @autoRoute
  @doc("List a {name}", Resource)
  @get
  @listsResource(Resource)
  @segmentOf(Resource)
  @armResourceList(Resource)
  List<
    Resource extends Foundations.SimpleResource,
    Parameters extends {} = {},
    Response extends {} = ArmResponse<ResourceListResult<Resource>>
  >(...ParentParameters, ...Parameters): Response | ErrorType;

  /**
   * A synchronous resource action.
   * @template Resource The resource being acted upon
   * @template Request The request model for the action
   * @template Response The response model for the action
   * @template Parameters Optional. Additional parameters after the path parameters
   * @template OptionalRequestBody Optional. Indicates whether the request body is optional
   */
  @doc("")
  @autoRoute
  @armResourceAction(Resource)
  @Private.enforceConstraint(Resource, Foundations.Resource)
  @post
  @returnsDoc("Azure operation completed successfully.")
  ActionSync<
    Resource extends Foundations.SimpleResource,
    Request extends TypeSpec.Reflection.Model | void,
    Response extends TypeSpec.Reflection.Model | void,
    Parameters extends {} = {},
    OptionalRequestBody extends valueof boolean = false
  >(
    ...ParentParameters,
    ...ResourceTypeParameter,
    ...Parameters,

    @doc("The content of the action request")
    @armBodyRoot(OptionalRequestBody)
    body: Request,
  ): Response | ErrorType;

  /**
   * A long-running resource action.
   * @template Resource The resource being acted upon
   * @template Request The request model for the action
   * @template LroHeaders Optional. Allows overriding the headers returned in the Accepted response
   * @template Response The response model for the action
   * @template Parameters Optional. Additional parameters after the path parameters
   * @template OptionalRequestBody Optional. Indicates whether the request body is optional
   */
  #suppress "@azure-tools/typespec-azure-core/no-response-body" "ARM"
  @doc("")
  @autoRoute
  @armResourceAction(Resource)
  @Private.enforceConstraint(Resource, Foundations.Resource)
  @post
  @returnsDoc("Azure operation completed successfully.")
  ActionAsync<
    Resource extends Foundations.SimpleResource,
    Request extends TypeSpec.Reflection.Model | void,
    Result extends TypeSpec.Reflection.Model | void,
    LroHeaders extends TypeSpec.Reflection.Model = ArmLroLocationHeader<
      Azure.Core.StatusMonitorPollingOptions<ArmOperationStatus>,
      Result,
      string
    > &
      Azure.Core.Foundations.RetryAfterHeader,
    Parameters extends {} = {},
    Response extends {} | void = ArmAcceptedLroResponse<
      "Resource operation accepted.",
      LroHeaders
    > | Result,
    OptionalRequestBody extends valueof boolean = false
  >(
    ...ParentParameters,
    ...ResourceTypeParameter,
    ...Parameters,

    @doc("The content of the action request")
    @armBodyRoot(OptionalRequestBody)
    body: Request,
  ): Response | ErrorType;
}
/**
 * @dev Get the provider namespace key-value pair
 * @template Resource Optional.  The resource to get the provider namespace for.
 */
model Provider<Resource extends {} = TenantActionScope> {
  ...ProviderNamespace<Resource>;
}

/**
 * A long-running resource update using a custom PATCH payload (Asynchronous)
 * @template Resource the resource being patched
 * @template PatchModel The input model for the PATCH request
 * @template BaseParameters Optional. Allows overriding the operation parameters
 * @template LroHeaders Optional.  Allows overriding the lro headers returned in the Accepted response
 * @template Parameters Optional. Additional parameters after the path parameters
 * @template Response Optional. The success response for the patch operation
 * @template Error Optional. The error response, if non-standard.
 * @template OptionalRequestBody Optional. Indicates whether the request body is optional
 */
@autoRoute
@doc("Update a {name}", Resource)
@armResourceUpdate(Resource)
@Private.enforceConstraint(Resource, Foundations.Resource)
@patch(#{ implicitOptionality: true })
op CustomPatchAsync<
  Resource extends Foundations.SimpleResource,
  PatchModel extends TypeSpec.Reflection.Model = TagsUpdateModel<Resource>,
  BaseParameters = DefaultBaseParameters<Resource>,
  LroHeaders extends TypeSpec.Reflection.Model = ArmLroLocationHeader<
    Azure.Core.StatusMonitorPollingOptions<ArmOperationStatus>,
    Resource,
    string
  > &
    Azure.Core.Foundations.RetryAfterHeader,
  Parameters extends {} = {},
  Response extends {} = ArmResponse<Resource> | ArmAcceptedLroResponse<
    "Resource update request accepted.",
    LroHeaders
  >,
  Error extends {} = ErrorResponse,
  OptionalRequestBody extends valueof boolean = false
> is UpdateOperation<
  ResourceInstanceParameters<Resource, BaseParameters> & Parameters,
  PatchModel,
  Response,
  Error,
  OptionalRequestBody
>;

/**
 * A resource update using a custom PATCH payload (synchronous)
 * @template Resource the resource being patched
 * @template PatchModel The input model for the PATCH request
 * @template BaseParameters Optional. Allows overriding the operation parameters
 * @template Parameters Optional. Additional parameters after the path parameters
 * @template Response Optional. The success response for the patch operation
 * @template Error Optional. The error response, if non-standard.
 * @template OptionalRequestBody Optional. Indicates whether the request body is optional
 */
@autoRoute
@doc("Update a {name}", Resource)
@armResourceUpdate(Resource)
@Private.enforceConstraint(Resource, Foundations.Resource)
@patch(#{ implicitOptionality: true })
op CustomPatchSync<
  Resource extends Foundations.SimpleResource,
  PatchModel extends TypeSpec.Reflection.Model = TagsUpdateModel<Resource>,
  BaseParameters = DefaultBaseParameters<Resource>,
  Parameters extends {} = {},
  Response extends {} = ArmResponse<Resource>,
  Error extends {} = ErrorResponse,
  OptionalRequestBody extends valueof boolean = false
> is UpdateOperation<
  ResourceInstanceParameters<Resource, BaseParameters> & Parameters,
  PatchModel,
  Response,
  Error,
  OptionalRequestBody
>;

/**
 * @dev The base template for Legacy Azure Resource Manager PATCH Operations.
 * @template HttpParameters The parameter object for the operation.
 * @template BodyParameter The body parameter
 * @template Response The response or union of responses for success.
 * @template ErrorResponse The error response.
 * @template OptionalRequestBody Optional. Indicates whether the request body is optional
 */
op UpdateOperation<
  HttpParameters extends {},
  BodyParameter extends {},
  Response extends {},
  ErrorResponse extends {},
  OptionalRequestBody extends valueof boolean = false
>(
  ...HttpParameters,

  @doc("The resource properties to be updated.")
  @armBodyRoot(OptionalRequestBody)
  properties: BodyParameter,
): Response | ErrorResponse;
