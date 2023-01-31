package handlers

import (
	"fmt"
	"mud/entities"
	"mud/services/tmap"
	"mud/services/variant"
)

func ParseWallVariant(t entities.Tile, rid, x, y, z int) string {
	var u, d, l, r bool
	var vid int

	_, err := fmt.Sscanf(t.Icon, "%d", &vid)
	if err != nil {
		return "ERR!"
	}

	tiles := tmap.GetSurroundingTiles(rid, x, y)
	for _, st := range tiles {
		if st.Z == z {
			if st.Y == y {
				if st.X == x-1 {
					l = true
				} else if st.X == x+1 {
					r = true
				}
			} else if st.X == x {
				if st.Y == y-1 {
					u = true
				} else if st.Y == y+1 {
					d = true
				}
			}
		}
	}

	var v entities.TileVariant
	var tname string = "V"
	if u {
		if l {
			if r {
				if d {
					tname = "Cross"
				} else {
					tname = "UT"
				}
			} else {
				if d {
					tname = "LT"
				} else {
					tname = "BRC"
				}
			}
		} else {
			if r {
				if d {
					tname = "RT"
				} else {
					tname = "BLC"
				}
			}
		}
	} else {
		tname = "H"
		if l {
			if r {
				if d {
					tname = "DT"
				}
			} else {
				if d {
					tname = "TRC"
				}
			}
		} else {
			if r {
				if d {
					tname = "TLC"
				}
			}
		}
	}

	vi := variant.CRUD.Retrieve(vid, tname)
	if vi == nil {
		return "ERR!"
	}

	v = vi.(entities.TileVariant)
	return v.Icon
}
