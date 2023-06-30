package ast

import (
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/ethereum-optimism/optimism/op-bindings/solc"
)

var remapTypeRe = regexp.MustCompile(`^(t_[\w_]+\([\w]+\))([\d]+)(_[\w]+)?$`)
var remapAstIdStorage = regexp.MustCompile(`(t_(struct|userDefinedValueType))\(([\w]+)\)([\d]+)_storage`)

// typeRemapping represents a mapping between an a type generated by solc
// and a canonicalized type. This is because solc inserts the ast id into
// certain types.
type typeRemapping struct {
	oldType string
	newType string
}

// CanonicalizeASTIDs canonicalizes AST IDs in storage layouts so that they
// don't cause unnecessary conflicts/diffs. The implementation is not
// particularly efficient, but is plenty fast enough for our purposes.
// It works in two passes:
//
//  1. First, it finds all AST IDs in storage and types, and builds a
//     map to replace them in the second pass.
//  2. The second pass performs the replacement.
//
// This function returns a copy of the passed-in storage layout. The
// inefficiency comes from replaceType, which performs a linear
// search of all replacements when performing substring matches of
// composite types.
func CanonicalizeASTIDs(in *solc.StorageLayout) *solc.StorageLayout {
	lastId := uint(1000)
	astIDRemappings := make(map[uint]uint)
	typeRemappings := make(map[string]string)

	for _, slot := range in.Storage {
		astIDRemappings[slot.AstId] = lastId
		lastId++
	}

	// Go map iteration order is random, so we need to sort
	// keys here in order to prevent non-determinism.
	var sortedOldTypes sort.StringSlice
	for oldType := range in.Types {
		sortedOldTypes = append(sortedOldTypes, oldType)
	}
	sortedOldTypes.Sort()

	seenTypes := make(map[string]bool)
	for _, oldType := range sortedOldTypes {
		if seenTypes[oldType] || oldType == "" {
			continue
		}

		matches := remapTypeRe.FindAllStringSubmatch(oldType, -1)
		if len(matches) == 0 {
			continue
		}

		// The storage types include the size when its a fixed size.
		// This is subject to breaking in the future if a type with
		// an ast id is added in a fixed storage type. We don't want
		// to skip a type with `_storage` in it if it has a subtype
		// with an ast id or it has an astid itself.
		skip := len(remapAstIdStorage.FindAllStringSubmatch(oldType, -1)) == 0
		if strings.Contains(oldType, "storage") && skip {
			continue
		}

		replaceAstID := matches[0][2]
		newType := strings.Replace(oldType, replaceAstID, strconv.Itoa(int(lastId)), 1)
		typeRemappings[oldType] = newType
		lastId++
		seenTypes[oldType] = true
	}

	outLayout := &solc.StorageLayout{
		Types: make(map[string]solc.StorageLayoutType),
	}
	for _, slot := range in.Storage {
		contract := slot.Contract

		// Normalize the name of the contract since absolute paths
		// are used when there are 2 contracts imported with the same
		// name
		if filepath.IsAbs(contract) {
			elements := strings.Split(contract, "/")
			if idx := slices.Index(elements, "optimism"); idx != -1 {
				contract = filepath.Join(elements[idx+1:]...)
			}
		}

		outLayout.Storage = append(outLayout.Storage, solc.StorageLayoutEntry{
			AstId:    astIDRemappings[slot.AstId],
			Contract: contract,
			Label:    slot.Label,
			Offset:   slot.Offset,
			Slot:     slot.Slot,
			Type:     replaceType(typeRemappings, slot.Type),
		})
	}

	for _, oldType := range sortedOldTypes {
		value := in.Types[oldType]
		newType := replaceType(typeRemappings, oldType)
		layout := solc.StorageLayoutType{
			Encoding:      value.Encoding,
			Label:         value.Label,
			NumberOfBytes: value.NumberOfBytes,
			Key:           replaceType(typeRemappings, value.Key),
			Value:         replaceType(typeRemappings, value.Value),
		}
		if value.Base != "" {
			layout.Base = replaceType(typeRemappings, value.Base)
		}
		outLayout.Types[newType] = layout

	}
	return outLayout
}

func replaceType(typeRemappings map[string]string, in string) string {
	if remap := typeRemappings[in]; remap != "" {
		return remap
	}

	// Track the number of matches
	matches := []typeRemapping{}
	for oldType, newType := range typeRemappings {
		if strings.Contains(in, oldType) {
			matches = append(matches, typeRemapping{oldType, newType})
		}
	}

	for _, match := range matches {
		in = strings.Replace(in, match.oldType, match.newType, 1)
	}

	return in
}
