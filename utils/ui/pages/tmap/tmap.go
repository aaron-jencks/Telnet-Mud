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
		r := rint.(entities.Room)

		portRangeX := utils.MAP_W - 2>>1
		portRangeY := utils.MAP_H - 2>>1
		smallWidth := utils.MAP_W-2 > r.Width
		smallHeight := utils.MAP_H-2 > r.Height
		nearLeft := roomX < portRangeX
		nearTop := roomY < portRangeY
		nearBottom := roomY >= r.Height-portRangeY
		nearRight := roomX >= r.Width-portRangeX

		// Setup view port
		trX := roomX - portRangeX
		if nearLeft {
			trX = 0
		}

		trY := roomY - portRangeY
		if nearTop {
			trY = 0
		}

		blX := roomX + portRangeX - 1
		if nearLeft {
			blX = utils.MAP_W - 2 - 1
		} else if nearRight {
			blX = r.Width - 1
		}

		blY := roomY + portRangeY - 1
		if nearTop {
			blY = utils.MAP_H - 2 - 1
		} else if nearBottom {
			blY = r.Height - 1
		}

		if smallWidth {
			diffX := utils.MAP_W - 2 - r.Width
			trX = -(diffX >> 1)
			blX = r.Width + diffX>>1
			if diffX%2 == 1 {
				blX++
			} else {
				blX--
			}
		}

		if smallHeight {
			diffY := utils.MAP_W - 2 - r.Width
			trY = -(diffY >> 1)
			blY = r.Width + diffY>>1
			if diffY%2 == 1 {
				blY++
			} else {
				blY--
			}
		}

		return trX, trY, blX, blY
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
	trx, try, blx, bly := GetMapPortCoords(p.Room, p.RoomX, p.RoomY)

	tiles := tmap.GetTilesForRegion(p.Room, trx, try, blx, bly)

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

	for _, mtile := range tiles {
		tilent := tile.CRUD.Retrieve(mtile.Tile).(entities.Tile)

		if mtile.Z > currentPort[mtile.Y][mtile.X].IconZ {
			// TODO add variant parsing here

			currentPort[mtile.Y][mtile.X].Icon = tilent.Icon
			currentPort[mtile.Y][mtile.X].FG = tilent.FG
			currentPort[mtile.Y][mtile.X].IconZ = mtile.Z
		}

		if tilent.BG > 0 && mtile.Z > currentPort[mtile.Y][mtile.X].BGZ {
			currentPort[mtile.Y][mtile.X].BG = tilent.BG
			currentPort[mtile.Y][mtile.X].BGZ = mtile.Z
		}
	}

	currentPort[p.RoomY][p.RoomX].Icon = utils.PLAYER_ICON
	currentPort[p.RoomY][p.RoomX].FG = utils.PLAYER_ICON_COLOR

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
