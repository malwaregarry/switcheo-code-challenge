package crude

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "crude/api/crude/crude"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: modulev1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod: "ResourceAll",
					Use:       "list-resource",
					Short:     "List all resource",
				},
				{
					RpcMethod:      "Resource",
					Use:            "show-resource [id]",
					Short:          "Shows a resource by id",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              modulev1.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "CreateResource",
					Use:            "create-resource [name] [category]",
					Short:          "Create resource",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "name"}, {ProtoField: "category"}},
				},
				{
					RpcMethod:      "UpdateResource",
					Use:            "update-resource [id] [name] [category]",
					Short:          "Update resource",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}, {ProtoField: "name"}, {ProtoField: "category"}},
				},
				{
					RpcMethod:      "DeleteResource",
					Use:            "delete-resource [id]",
					Short:          "Delete resource",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
