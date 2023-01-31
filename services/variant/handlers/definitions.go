package handlers

import "mud/entities"

// Defines a variant handler for a given variant
// Takes the given tile, the room id, and the tile's location
// The function returns the corresponding variant icon string for that tile
type VariantHandler func(t entities.Tile, rId, x, y, z int) string
