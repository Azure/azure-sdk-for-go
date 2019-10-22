package azidentity

// THIS CREDENTIAL BUILDS ON TOP OF MSAL
// import (
// 	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
// )

// type PublicClientApplication interface {
// 	IsSystemWebViewAvailable bool

// }
// //  UsernamePasswordCredential enables authentication to Azure Active Directory using a user's  username and password. If the user has MFA enabled this
// //  credential will fail to get a token and return an AuthenticationFailedError. Also, this credential requires a high degree of
// //  trust and is not recommended outside of prototyping when more secure credentials can be used.
// type UsernamePasswordCredential struct { // implements TokenCredential
// 	// pubApp   IPublicClientApplication
// 	pipeline azcore.Pipeline
// 	options  IdentityClientOptions
// 	username string
// 	password string
// }

// // NewUsernamePasswordCredential creates an instance of the UsernamePasswordCredential with the details needed to authenticate against Azure Active Directory with a simple username
// // and password.
// // Username
// // Password: The user account's user name, UPN.
// // clientID: The client (application) ID of an App Registration in the tenant.
// // tenantID: The Azure Active Directory tenant (directory) ID or name.
// // Options: The client options for the newly created UsernamePasswordCredential.
// func NewUsernamePasswordCredential(username string, password string, clientID string, tenantID string, options IdentityClientOptions) UsernamePasswordCredential {
// 	// CP: How do we return errors for these cases?
// 	if username == "" {
// 		return UsernamePasswordCredential{} //CP fix this create error type with msg param
// 	}
// 	if password == "" {
// 		return UsernamePasswordCredential{}
// 	}
// 	// if options == nil {
// 	// 	options = NewIdentityClientOptions()
// 	// }
// 	pipeline := NewDefaultPipeline(azcore.PipelineOptions{}, nil)
// 	// CP: changed new HttpPipelineClient to NewHttpPipelineClient, probably need to call this on an instance of that object
// 	// CP: I need to implement HttpPipelineClient?
// 	pubApp := PublicClientApplicationBuilder.Create(clientId).WithHttpClientFactory(NewHttpPipelineClientFactory(pipeline)).WithTenantId(tenantId).Build()
// 	return UsernameUsernamePasswordCredential{pubApp: pubApp, pipeline: pipeline, options: options, username: username, password: password}
// }

// // Obtains a token for a user account, authenticating them using the given username and password.  Note: This will fail with
// // an AuthenticationFailedError if the specified user acound has MFA enabled.
// // Scopes: The list of scopes for which the token will have access.
// // Returns AccessToken which can be used to authenticate service client calls.
// // CP: Do we want this to be async?
// func (c UsernamePasswordCredential) GetToken(ctx context.Context, scopes []string) (AccessToken, error) {
// 	// var scope DiagnosticScope
// 	// scope = pipeline.Diagnostics.CreateScope("Azure.Identity.UsernamePasswordCredential.GetToken")
// 	// scope.Start()

// 	var result AuthenticationResult
// 	// CP: TO-DO this needs to check for an error and  return an error is the AuthenticationResult was not created correctly
// 	result = pubApp.AcquireTokenByUsernamePassword(scopes, _username, _password).ExecuteAsync(cancellationToken).ConfigureAwait(false)
// 	return NewAccessToken(result.AccessToken, result.ExpiresOn)
// }
