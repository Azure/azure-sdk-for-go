import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/automanage/armautomanage"
)

var (
	subId       = "<sub id>"
	rg          = "resourceGroupName"
	vm          = "vmName"
	profileName = "custom-profile"

	// ------------------------- SET UP AUTOMANAGE CLIENT--------------------------------
	credential, _           = azidentity.NewDefaultAzureCredential(nil)
	configProfilesClient, _ = armautomanage.NewConfigurationProfilesClient(subId, credential, nil)
	assignmentClient, _     = armautomanage.NewConfigurationProfileAssignmentsClient(subId, credential, nil)
)

func main() {
	// ------------------------GET PROFILE------------------------------------------------
	profile, _ := configProfilesClient.Get(context.Background(), profileName, rg, nil)
	data, _ := json.MarshalIndent(profile, "", "   ")

	fmt.Println(string(data))

	// ------------------------GET ALL PROFILES IN RESOURCE GROUP-------------------------
	profiles, _ := configProfilesClient.NewListByResourceGroupPager(rg, nil).NextPage(context.Background())
	data, _ := json.MarshalIndent(profiles, "", "   ")

	fmt.Println(string(data))

	// ------------------------GET ALL PROFILES IN SUBSCRIPTION---------------------------
	profiles, _ := configProfilesClient.NewListBySubscriptionPager(nil).NextPage(context.Background())
	data, _ := json.MarshalIndent(profiles, "", "   ")

	fmt.Println(string(data))

	// ------------------------CREATE OR UPDATE CUSTOM PROFILE----------------------------

	// to update, provide a value for all properties as if you were creating a configuration profile (ID, Name, Type, Location, Properties, Tags)
	config := make(map[string]interface{})
	config["Antimalware/Enable"] = true
	config["Antimalware/Exclusions/Paths"] = ""
	config["Antimalware/Exclusions/Extensions"] = ""
	config["Antimalware/Exclusions/Processes"] = ""
	config["Antimalware/EnableRealTimeProtection"] = false
	config["Antimalware/RunScheduledScan"] = true
	config["Antimalware/ScanType"] = "Quick"
	config["Antimalware/ScanDay"] = 7
	config["Antimalware/ScanTimeInMinutes"] = 120
	config["Backup/Enable"] = true
	config["Backup/PolicyName"] = "dailyBackupPolicy"
	config["Backup/TimeZone"] = "UTC"
	config["Backup/InstantRpRetentionRangeInDays"] = 2
	config["Backup/SchedulePolicy/ScheduleRunFrequency"] = "Daily"
	config["Backup/SchedulePolicy/ScheduleRunTimes"] = []string{"2022-07-27T12: 00: 00Z"}
	config["Backup/SchedulePolicy/SchedulePolicyType"] = "SimpleSchedulePolicy"
	config["Backup/RetentionPolicy/RetentionPolicyType"] = "LongTermRetentionPolicy"
	config["Backup/RetentionPolicy/DailySchedule/RetentionTimes"] = []string{"2022-07-27T12: 00: 00Z"}
	config["Backup/RetentionPolicy/DailySchedule/RetentionDuration/Count"] = 180
	config["Backup/RetentionPolicy/DailySchedule/RetentionDuration/DurationType"] = "Days"
	config["WindowsAdminCenter/Enable"] = false
	config["VMInsights/Enable"] = true
	config["AzureSecurityCenter/Enable"] = true
	config["UpdateManagement/Enable"] = true
	config["ChangeTrackingAndInventory/Enable"] = true
	config["GuestConfiguration/Enable"] = true
	config["AutomationAccount/Enable"] = true
	config["LogAnalytics/Enable"] = true
	config["BootDiagnostics/Enable"] = true
	configJson, _ := json.MarshalIndent(config, "", "   ")
	fmt.Println(string(configJson))

	properties := armautomanage.ConfigurationProfileProperties{
		Configuration: config,
	}

	id := fmt.Sprintf("/subscriptions/%v/resourceGroups/%v/providers/Microsoft.Automanage/configurationProfiles/%v", subId, rg, profileName)
	resourceType := "Microsoft.Automanage/configurationProfiles"
	location := "eastus"
	env := "prod"

	tags := make(map[string]*string)
	tags["environment"] = &env

	newProfile := armautomanage.ConfigurationProfile{
		ID:         &id,
		Name:       &profileName,
		Type:       &resourceType,
		Location:   &location,
		Properties: &properties,
		Tags:       tags,
	}

	configProfilesClient.CreateOrUpdate(context.Background(), profileName, rg, newProfile, nil)

	// ------------------------ DELETE PROFILE  ------------------------------------------
	configProfilesClient.Delete(context.Background(), rg, profileName, nil)

	// ------------------------GET ASSIGNMENT---------------------------------------------
	assignment, _ := assignmentClient.Get(context.Background(), rg, "default", vm, nil)
	data, _ := json.MarshalIndent(assignment, "", "   ")
	fmt.Println(string(data))

	// ------------------------GET ALL ASSIGNMENTS IN A SUBSCRIPTION ---------------------
	assignments, _ := assignmentClient.NewListBySubscriptionPager(nil).NextPage(context.Background())
	data, _ := json.MarshalIndent(assignments, "", "   ")
	fmt.Println(string(data))

	// ------------------------CREATE BEST PRACTICES PRODUCTION PROFILE ASSIGNMENT--------
	vmId := fmt.Sprintf("/subscriptions/%v/resourceGroups/%v/providers/Microsoft.Compute/virtualMachines/%v", subId, rg, vm)
	configProfileId := "/providers/Microsoft.Automanage/bestPractices/AzureBestPracticesProduction"

	properties := armautomanage.ConfigurationProfileAssignmentProperties{
		ConfigurationProfile: &configProfileId,
		TargetID:             &vmId,
	}

	id := fmt.Sprintf("/subscriptions/%v/resourceGroups/%v/providers/Microsoft.Compute/virtualMachines/%v/providers/Microsoft.Automanage/AutomanageAssignments/default", subId, rg, vm)
	name := "default"
	assignment := armautomanage.ConfigurationProfileAssignment{
		ID:         &id,
		Name:       &name,
		Properties: &properties,
	}

	assignmentClient.CreateOrUpdate(context.Background(), "default", rg, vm, assignment, nil)

	// ------------------------CREATE CUSTOM PROFILE ASSIGNMENT---------------------------
	vmId := fmt.Sprintf("/subscriptions/%v/resourceGroups/%v/providers/Microsoft.Compute/virtualMachines/%v", subId, rg, vm)
	configProfileId := fmt.Sprintf("/subscriptions/%v/resourceGroups/%v/providers/Microsoft.Automanage/configurationProfiles/%v", subId, rg, profileName)

	properties := armautomanage.ConfigurationProfileAssignmentProperties{
		ConfigurationProfile: &configProfileId,
		TargetID:             &vmId,
	}

	id := fmt.Sprintf("/subscriptions/%v/resourceGroups/%v/providers/Microsoft.Compute/virtualMachines/%v/providers/Microsoft.Automanage/AutomanageAssignments/default", subId, rg, vm)
	name := "default"
	assignment := armautomanage.ConfigurationProfileAssignment{
		ID:         &id,
		Name:       &name,
		Properties: &properties,
	}

	assignmentClient.CreateOrUpdate(context.Background(), "default", rg, vm, assignment, nil)
}
