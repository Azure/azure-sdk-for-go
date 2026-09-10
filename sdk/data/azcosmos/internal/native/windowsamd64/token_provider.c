#include "../../../azurecosmosdriver.h"

extern int32_t goCosmosTokenProviderGet(
    intptr_t user_data,
    const cosmos_token_request_t *request);
extern void goCosmosTokenProviderFree(intptr_t user_data);

cosmos_token_provider_t cosmos_go_token_provider(void)
{
    cosmos_token_provider_t provider = {
        .get_token = goCosmosTokenProviderGet,
        .user_data_free = goCosmosTokenProviderFree,
    };
    return provider;
}
