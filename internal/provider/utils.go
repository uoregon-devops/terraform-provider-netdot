package provider

import (
	"context"
	"fmt"
	"terraform-provider-netdot/internal/netdot/models"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func isPopulated(value attr.Value) bool {
	return !value.IsNull() && !value.IsUnknown()
}

func autoNullString(s string) types.String {
	if s == "" {
		return types.StringNull()
	}
	return types.StringValue(s)
}

func autoNullInt64(i int64) types.Int64 {
	if i == 0 {
		return types.Int64Null()
	}
	return types.Int64Value(i)
}

// autonull xlink string
func autoNullXlinkString(xlink models.Xlink, s string) types.String {
	if xlink.ID == 0 {
		return types.StringNull()
	}
	return types.StringValue(s)
}

// autonull xlink int64
func autoNullXlinkInt64(xlink models.Xlink, i int64) types.Int64 {
	if xlink.ID == 0 {
		return types.Int64Null()
	}
	return types.Int64Value(i)
}

type fieldReactivePlanModifier struct {
	mappedFieldPath   string
	mappedFieldBuffer attr.Value
}

func (pm fieldReactivePlanModifier) Description(ctx context.Context) string {
	return fmt.Sprintf("Reactive string plan modifier for field %s", pm.mappedFieldPath)
}

func (pm fieldReactivePlanModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Reactive string plan modifier for field %s", pm.mappedFieldPath)
}

func (pm fieldReactivePlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	stateValue := pm.mappedFieldBuffer
	req.State.GetAttribute(ctx, path.Root(pm.mappedFieldPath), &stateValue)

	configValue := pm.mappedFieldBuffer
	req.Config.GetAttribute(ctx, path.Root(pm.mappedFieldPath), &configValue)

	// Check if the field you care about has changed
	if !stateValue.Equal(configValue) {
		// Mark this field as unknown so it will be recalculated
		resp.PlanValue = types.StringUnknown()
	} else {
		// Otherwise, use the existing value
		resp.PlanValue = req.StateValue
	}
}

func (pm fieldReactivePlanModifier) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	stateValue := pm.mappedFieldBuffer
	req.State.GetAttribute(ctx, path.Root(pm.mappedFieldPath), &stateValue)

	planValue := pm.mappedFieldBuffer
	req.Plan.GetAttribute(ctx, path.Root(pm.mappedFieldPath), &planValue)

	// Check if the field you care about has changed
	if !stateValue.Equal(planValue) {
		// Mark this field as unknown so it will be recalculated
		resp.PlanValue = types.Int64Unknown()
	} else {
		// Otherwise, use the existing value
		resp.PlanValue = req.StateValue
	}
}

func CreateFieldReactiveStringPlanModifier(mappedFieldPath string, attributeType attr.Value) planmodifier.String {
	return fieldReactivePlanModifier{
		mappedFieldPath:   mappedFieldPath,
		mappedFieldBuffer: attributeType,
	}
}

func CreateFieldReactiveInt64PlanModifier(mappedFieldPath string, attributeType attr.Value) planmodifier.Int64 {
	return fieldReactivePlanModifier{
		mappedFieldPath:   mappedFieldPath,
		mappedFieldBuffer: attributeType,
	}
}
