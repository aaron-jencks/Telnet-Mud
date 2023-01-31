package tmap

import (
	"fmt"
	"mud/entities"
	"mud/services/room"
	"mud/services/tile"
	"mud/services/tmap"
	"mud/utils"
	"mud/utils/ui"
	"mud/utils/ui/gui"
	"strings"
)

func GetMapPortCoords(roomId, roomX, roomY int) (int, int, int, int) {
	rint := room.CRUD.Retrieve(roomId)
	if rint != nil {
		r := rint.(room.ExpandedRoom)

		// The view distance of the player
		portRangeX := (utils.MAP_W - 2) >> 1
		portRangeY := (utils.MAP_H - 2) >> 1

		smallWidth := utils.MAP_W-2 > r.Width
		smallHeight := utils.MAP_H-2 > r.Height
		if smallWidth {
			portRangeX = r.Width >> 1
		}
		if smallHeight {
			portRangeY = r.Height >> 1
		}

		nearLeft := roomX < portRangeX
		nearTop := roomY < portRangeY
		nearBottom := roomY >= r.Height-portRangeY
		nearRight := roomX >= r.Width-portRangeX

		// Setup view port

		var tlX, tlY, brX, brY int

		// X
		if nearLeft || smallWidth {
			tlX = 0
			if smallWidth {
				brX = r.Width - 1
			} else {
				brX = utils.MAP_W - 3
			}
		} else if nearRight {
			tlX = r.Width - 1 - portRangeX
			brX = r.Width - 1
		} else {
			tlX = roomX - portRangeX
			brX = roomX + portRangeX - 1
		}

		// Y
		if nearTop || smallHeight {
			tlY = 0
			if smallHeight {
				brY = r.Height - 1
			} else {
				brY = utils.MAP_H - 3
			}
		} else if nearBottom {
			tlY = r.Height - 1 - portRangeY
			brY = r.Height - 1
		} else {
			tlY = roomY - portRangeY
			brY = roomY + portRangeY - 1
		}

		return tlX, tlY, brX, brY
	}

	return 0, 0, utils.MAP_W - 2 - 1, utils.MAP_H - 2 - 1
}

type TileInfo struct {
	Icon  string
	FG    int
	BG    int
	IconZ int
	BGZ   int
}

func (ti TileInfo) ToString() string {
	return ui.AddBackground(ti.BG, ui.CSI(fmt.Sprint(ti.FG), "m")+ti.Icon)
}

func GetMapWindow(p entities.Player) string {
	tlx, tly, brx, bry := GetMapPortCoords(p.Room, p.RoomX, p.RoomY)

	var xoffset, yoffset int
	if utils.MAP_W-2 > (brx - tlx) {
		xoffset = (utils.MAP_W - 2 - (brx - tlx)) >> 1
	}
	if utils.MAP_H-2 > (bry - tly) {
		yoffset = (utils.MAP_H - 2 - (bry - tly)) >> 1
	}

	tiles := tmap.GetTilesForRegion(p.Room, tlx, tly, brx, bry)

	var currentPort [][]TileInfo = make([][]TileInfo, utils.MAP_H-2)
	for row := 0; row < utils.MAP_H-2; row++ {
		currentPort[row] = make([]TileInfo, utils.MAP_W-2)
		for col := 0; col < utils.MAP_W-2; col++ {
			currentPort[row][col] = TileInfo{
				Icon:  utils.DEFAULT_MAP_BACKGROUND,
				FG:    utils.DEFAULT_MAP_BACKGROUND_FG_COLOR,
				BG:    utils.DEFAULT_MAP_BACKGROUND_BG_COLOR,
				IconZ: -1,
				BGZ:   -1,
			}
		}
	}

	rent := room.CRUD.Retrieve(p.Room)
	if rent != nil {
		rentst := rent.(room.ExpandedRoom)

		defHasBG := rentst.BackgroundTile.BG > 0

		for trow := 0; trow < bry-tly; trow++ {
			yindex := trow - tly + yoffset
			for tcol := 0; tcol < brx-tlx; tcol++ {
				xindex := tcol - tlx + xoffset
				currentPort[yindex][xindex] = TileInfo{
					Icon:  rentst.BackgroundTile.Icon,
					FG:    rentst.BackgroundTile.FG,
					BG:    utils.DEFAULT_MAP_BACKGROUND_BG_COLOR,
					IconZ: -1,
					BGZ:   -1,
				}

				if defHasBG {
					currentPort[yindex][xindex].BG = rentst.BackgroundTile.BG
				}
			}
		}
	}

	for _, mtile := range tiles {
		tilent := tile.CRUD.Retrieve(mtile.Tile).(entities.Tile)

		yindex := mtile.Y - tly + yoffset
		xindex := mtile.X - tlx + xoffset

		if mtile.Z > currentPort[yindex][xindex].IconZ {
			// TODO add variant parsing here

			currentPort[yindex][xindex].Icon = tilent.Icon
			currentPort[yindex][xindex].FG = tilent.FG
			currentPort[yindex][xindex].IconZ = mtile.Z
		}

		if tilent.BG > 0 && mtile.Z > currentPort[yindex][xindex].BGZ {
			currentPort[yindex][xindex].BG = tilent.BG
			currentPort[yindex][xindex].BGZ = mtile.Z
		}
	}

	pyindex := p.RoomY - tly + yoffset
	pxindex := p.RoomX - tlx + xoffset

	currentPort[pyindex][pxindex].Icon = utils.PLAYER_ICON
	currentPort[pyindex][pxindex].FG = utils.PLAYER_ICON_COLOR

	var rows []string = make([]string, len(currentPort))
	for ri, row := range currentPort {
		var stringRow []string = make([]string, len(row))
		for ci, col := range row {
			stringRow[ci] = col.ToString()
		}
		rows[ri] = strings.Join(stringRow, "")
	}

	return gui.SizedBoxText(strings.Join(rows, "\n"), utils.MAP_H, utils.MAP_W)
}
