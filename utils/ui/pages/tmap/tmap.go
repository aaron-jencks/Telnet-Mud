package tmap

import (
	"mud/entities"
	"mud/services/player"
	"mud/services/room"
	"mud/services/tile"
	"mud/services/tmap"
	"mud/utils"
	"mud/utils/ui/gui"
	"net"
	"strings"
)

func GetMapPortCoords(p entities.Player) (int, int, int, int) {
	rint := room.CRUD.Retrieve(p.Room)
	if rint != nil {
		r := rint.(entities.Room)

		portRangeX := utils.MAP_W - 2>>1
		portRangeY := utils.MAP_H - 2>>1
		smallWidth := utils.MAP_W-2 > r.Width
		smallHeight := utils.MAP_H-2 > r.Height
		nearLeft := p.RoomX < portRangeX
		nearTop := p.RoomY < portRangeY
		nearBottom := p.RoomY >= r.Height-portRangeY
		nearRight := p.RoomX >= r.Width-portRangeX

		// Setup view port
		trX := p.RoomX - portRangeX
		if nearLeft {
			trX = 0
		}

		trY := p.RoomY - portRangeY
		if nearTop {
			trY = 0
		}

		blX := p.RoomX + portRangeX - 1
		if nearLeft {
			blX = utils.MAP_W - 2 - 1
		} else if nearRight {
			blX = r.Width - 1
		}

		blY := p.RoomY + portRangeY - 1
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

func GetMapWindow(conn net.Conn) string {
	p := player.CRUD.Retrieve(player.PlayerConnectionMap[conn]).(entities.Player)

	trx, try, blx, bly := GetMapPortCoords(p)

	tiles := tmap.GetTilesForRegion(p.Room, trx, try, blx, bly)

	var currentPort [][]string = make([][]string, utils.MAP_H-2)
	var currentPortLevel [][]int = make([][]int, utils.MAP_H-2)
	for row := 0; row < utils.MAP_H-2; row++ {
		currentPort[row] = make([]string, utils.MAP_W-2)
		currentPortLevel[row] = make([]int, utils.MAP_W-2)
		for col := 0; col < utils.MAP_W-2; col++ {
			currentPort[row][col] = utils.DEFAULT_MAP_BACKGROUND
			currentPortLevel[row][col] = -1
		}
	}

	for _, mtile := range tiles {
		if mtile.Z > currentPortLevel[mtile.Y][mtile.X] {
			tilent := tile.CRUD.Retrieve(mtile.Tile).(entities.Tile)

			// TODO add variant parsing here

			currentPort[mtile.Y][mtile.X] = tilent.Icon
			currentPortLevel[mtile.Y][mtile.X] = mtile.Z
		}
	}

	var rows []string = make([]string, len(currentPort))
	for ri, row := range currentPort {
		rows[ri] = strings.Join(row, "")
	}

	return gui.SizedBoxText(strings.Join(rows, "\n"), utils.MAP_H, utils.MAP_W)
}
