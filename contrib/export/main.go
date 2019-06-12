package main

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	extypes "github.com/cosmos/cosmos-sdk/contrib/export/types"
	"github.com/cosmos/cosmos-sdk/contrib/export/v036"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/tendermint/tendermint/types"
)

// TODO accept those as parameters once we have more than one release available
var source = "034"
var target = "036"

var migrationMap = extypes.MigrationMap{
	"036": v036.Migrate,
}

func main() {
	// This function should be modularized, for now we read and dump genesis committed to git,
	// to simplify the creation of a CCI script that tests three different things:
	// - Reading old types,
	// - converting to new ones
	// - reading the committed new genesis and see if all works

	// We could add an invariant test for genesis, conversion should be identical to exporting the target genesis once imported

	cdc := codec.New()
	// Dump here an example genesis for CCI to test import from old types, and export to new ones
	genDoc, err := types.GenesisDocFromFile(fmt.Sprintf("./contrib/export/genesis_%s.json", source))
	if err != nil {
		panic(err)
	}
	var initialState extypes.AppMap
	cdc.MustUnmarshalJSON(genDoc.AppState, &initialState)

	if migrationMap[target] == nil {
		panic(fmt.Sprintf("Missing migration function for version %s", target))
	}
	newGenState := migrationMap[target](initialState, cdc)

	genDoc.AppState = cdc.MustMarshalJSON(newGenState)
	// Keep dumping updated genesis to test import of a new genesis directly
	if err = genutil.ExportGenesisFile(genDoc, fmt.Sprintf("./contrib/export/genesis_%s.json", target)); err != nil {
		panic(err)
	}
}
