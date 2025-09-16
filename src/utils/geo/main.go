package geo

import (
	"fmt"
	"math"
	"openstreetmap-go/src/repository/model"
	"sort"
	"strconv"
	"strings"
)

func Lon2x(lon float64) int {
	return int(math.Round((lon + 180.0) * 65535.0 / 360.0))
}

func Lat2y(lat float64) int {
	return int(math.Round((lat + 90.0) * 65535.0 / 180.0))
}

func XY2Tile(x, y int) uint32 {
	tile := 0
	for i := 15; i >= 0; i-- {
		tile = (tile << 1) | ((x >> i) & 1)
		tile = (tile << 1) | ((y >> i) & 1)
	}
	return uint32(tile)
}

func TileIdsInBoundingBox(minX, minY, maxX, maxY int) []int {
	tiles := []int{}
	for i := minX; i <= maxX; i++ {
		for j := minY; j <= maxY; j++ {
			tiles = append(tiles, int(XY2Tile(i, j)))
		}
	}
	sort.Slice(tiles, func(i, j int) bool { return tiles[i] < tiles[j] })
	return tiles
}

func GroupTiles(tileList []int) ([][]int, []int) {
	var ranges [][]int
	var singles []int

	if len(tileList) == 0 {
		return ranges, singles
	}

	start := tileList[0]
	end := start

	for i := 1; i < len(tileList); i++ {
		tile := tileList[i]
		if tile == end+1 {
			end = tile
		} else {
			if start == end {
				singles = append(singles, start)
			} else {
				ranges = append(ranges, []int{start, end})
			}
			start = tile
			end = tile
		}
	}

	if start == end {
		singles = append(singles, start)
	} else {
		ranges = append(ranges, []int{start, end})
	}

	return ranges, singles
}

func BuildTileAreaSQLClause(bbox model.BBox, prefix string) string {

	minX := Lon2x(bbox.MinLon)
	minY := Lat2y(bbox.MinLat)
	maxX := Lon2x(bbox.MaxLon)
	maxY := Lat2y(bbox.MaxLat)

	tileList := TileIdsInBoundingBox(minX, minY, maxX, maxY)
	ranges, singles := GroupTiles(tileList)

	sqlParts := []string{}
	for _, r := range ranges {
		start, end := r[0], r[1]
		sqlParts = append(sqlParts, fmt.Sprintf("%stile BETWEEN %d AND %d", prefix, start, end))
	}

	if len(singles) > 0 {
		singleStrs := make([]string, len(singles))
		for i, val := range singles {
			singleStrs[i] = strconv.Itoa(val)
		}
		sqlParts = append(sqlParts, fmt.Sprintf("%stile IN (%s)", prefix, strings.Join(singleStrs, ",")))
	}
	return fmt.Sprintf("( %s )", strings.Join(sqlParts, " OR "))
}
