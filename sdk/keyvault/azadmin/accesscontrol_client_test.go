//go:build go1.18
// +build go1.18

package azadmin_test

/*func TestCreateRoleAssignment(t *testing.T) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	require.NoError(t, err)

	client, err := azadmin.NewAccessControlClient("", cred, nil)
	require.NoError(t, err)
	_ = client

	res3, err := client.CreateRoleAssignment(context.Background(), "/", "testname", azadmin.RoleAssignmentCreateParameters{}, nil)
	require.NoError(t, err)
	_ = res3

}*/
