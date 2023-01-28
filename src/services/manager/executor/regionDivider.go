package executor

import (
	"context"

	"github.com/badochov/distributed-shortest-path/src/libs/db"
)

type regionDivider struct {
	db         db.DB
	generation generation
}

func (rd *regionDivider) doDivideIntoRegions(ctx context.Context, coordsBetween db.CoordBounds,
	count int64, minRegionId, maxRegionId db.RegionId, vertical bool) error {
	if minRegionId == maxRegionId {
		_, err := rd.db.SetRegion(ctx, coordsBetween, minRegionId, rd.generation)
		return err
	}
	if vertical {
		left, right := coordsBetween.Longitude.Min, coordsBetween.Longitude.Max
		var mid float64
		var leftPart, rightPart, midPart int64
		var err error
		for {
			mid = (left + right) / 2
			midPart, err = rd.db.GetVertexCountOnVerticalSegment(ctx, coordsBetween.Latitude, mid, rd.generation)
			if err != nil {
				return err
			}
			leftPart, err = rd.db.GetVertexCount(ctx,
				db.CoordBounds{
					Latitude:  coordsBetween.Latitude,
					Longitude: db.MinMax{Min: coordsBetween.Longitude.Min, Max: mid},
				},
				rd.generation)
			if err != nil {
				return err
			}
			rightPart = count - leftPart // including midPart
			if leftPart > rightPart {
				right = mid
			} else if leftPart == rightPart || leftPart+midPart > rightPart-midPart {
				break
			} else {
				left = mid
			}
		}
		err = rd.doDivideIntoRegions(ctx,
			db.CoordBounds{
				Latitude:  coordsBetween.Latitude,
				Longitude: db.MinMax{Min: coordsBetween.Longitude.Min, Max: mid},
			},
			leftPart, minRegionId, (minRegionId+maxRegionId)/2, false)
		if err != nil {
			return err
		}
		return rd.doDivideIntoRegions(ctx,
			db.CoordBounds{
				Latitude:  coordsBetween.Latitude,
				Longitude: db.MinMax{Min: mid, Max: coordsBetween.Longitude.Max}},
			rightPart, (minRegionId+maxRegionId)/2+1, maxRegionId, false)
	} else {
		down, up := coordsBetween.Latitude.Min, coordsBetween.Latitude.Max
		var mid float64
		var downPart, upPart, midPart int64
		var err error
		for {
			mid = (down + up) / 2
			midPart, err = rd.db.GetVertexCountOnHorizontalSegment(ctx, mid, coordsBetween.Longitude, rd.generation)
			if err != nil {
				return err
			}
			downPart, err = rd.db.GetVertexCount(ctx,
				db.CoordBounds{
					Latitude:  db.MinMax{Min: coordsBetween.Latitude.Min, Max: mid},
					Longitude: coordsBetween.Longitude,
				},
				rd.generation)
			if err != nil {
				return err
			}
			upPart = count - downPart // including midPart
			if downPart > upPart {
				up = mid
			} else if downPart == upPart || downPart+midPart > upPart-midPart {
				break
			} else {
				down = mid
			}
		}
		err = rd.doDivideIntoRegions(ctx,
			db.CoordBounds{
				Latitude:  db.MinMax{Min: coordsBetween.Latitude.Min, Max: mid},
				Longitude: coordsBetween.Longitude,
			},
			downPart, minRegionId, (minRegionId+maxRegionId)/2, true)
		if err != nil {
			return err
		}
		return rd.doDivideIntoRegions(ctx,
			db.CoordBounds{
				Latitude:  db.MinMax{Min: mid, Max: coordsBetween.Latitude.Max},
				Longitude: coordsBetween.Longitude,
			},
			upPart, (minRegionId+maxRegionId)/2+1, maxRegionId, true)
	}
}

func (rd *regionDivider) divideIntoRegions(ctx context.Context, numRegions int) error {
	bounds := db.CoordBounds{
		Latitude:  db.MinMax{Min: -90, Max: 91},
		Longitude: db.MinMax{Min: -180, Max: 181},
	}
	count, err := rd.db.GetVertexCount(ctx, bounds, rd.generation)
	if err != nil {
		return err
	}
	return rd.doDivideIntoRegions(ctx, bounds, count, 0, regionId(numRegions)-1, true)
}
