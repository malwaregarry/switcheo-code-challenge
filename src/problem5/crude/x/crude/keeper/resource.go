package keeper

import (
	"context"
	"encoding/binary"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"crude/x/crude/types"
	"github.com/cosmos/cosmos-sdk/runtime"
)

// GetResourceCount get the total number of resource
func (k Keeper) GetResourceCount(ctx context.Context) uint64 {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	byteKey := types.KeyPrefix(types.ResourceCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	return binary.BigEndian.Uint64(bz)
}

// SetResourceCount set the total number of resource
func (k Keeper) SetResourceCount(ctx context.Context, count uint64) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	byteKey := types.KeyPrefix(types.ResourceCountKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(byteKey, bz)
}

// AppendResource appends a resource in the store with a new id and update the count
func (k Keeper) AppendResource(
	ctx context.Context,
	resource types.Resource,
) uint64 {
	// Create the resource
	count := k.GetResourceCount(ctx)

	// Set the ID of the appended value
	resource.Id = count

	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ResourceKey))
	appendedValue := k.cdc.MustMarshal(&resource)
	store.Set(GetResourceIDBytes(resource.Id), appendedValue)

	// Update resource count
	k.SetResourceCount(ctx, count+1)

	return count
}

// SetResource set a specific resource in the store
func (k Keeper) SetResource(ctx context.Context, resource types.Resource) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ResourceKey))
	b := k.cdc.MustMarshal(&resource)
	store.Set(GetResourceIDBytes(resource.Id), b)
}

// GetResource returns a resource from its id
func (k Keeper) GetResource(ctx context.Context, id uint64) (val types.Resource, found bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ResourceKey))
	b := store.Get(GetResourceIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveResource removes a resource from the store
func (k Keeper) RemoveResource(ctx context.Context, id uint64) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ResourceKey))
	store.Delete(GetResourceIDBytes(id))
}

// GetAllResource returns all resource
func (k Keeper) GetAllResource(ctx context.Context) (list []types.Resource) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ResourceKey))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Resource
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetResourceIDBytes returns the byte representation of the ID
func GetResourceIDBytes(id uint64) []byte {
	bz := types.KeyPrefix(types.ResourceKey)
	bz = append(bz, []byte("/")...)
	bz = binary.BigEndian.AppendUint64(bz, id)
	return bz
}
