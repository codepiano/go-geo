package geo

import (
	"math"

	"github.com/flywave/go-geo/gcj02"
	vec2d "github.com/flywave/go3d/float64/vec2"
)

// GCJ02WebMercatorProj is spherical Web Mercator over GCJ-02 longitude/latitude,
// as used by AMap XYZ tiles. GCJ02:3857 is a custom identifier, not an EPSG code.
// The existing GCJ02MC projection retains its legacy polynomial behavior.
type GCJ02WebMercatorProj struct{}

func (p *GCJ02WebMercatorProj) TransformTo(o Proj, points []vec2d.T) []vec2d.T {
	if p.Eq(o) {
		return points
	}
	geographic := make([]vec2d.T, len(points))
	for i, point := range points {
		geographic[i] = vec2d.T{point[0] / 6378137 * 180 / math.Pi,
			(2*math.Atan(math.Exp(point[1]/6378137)) - math.Pi/2) * 180 / math.Pi}
	}
	if _, ok := o.(*BD09Proj); ok {
		for i, point := range geographic {
			lat, lon := gcj02.GCJ02toBD09(point[1], point[0])
			geographic[i] = vec2d.T{lon, lat}
		}
		return geographic
	}
	return newGCJ02Proj(true).TransformTo(o, geographic)
}

func (p *GCJ02WebMercatorProj) transformFromGCJ02(points []vec2d.T) []vec2d.T {
	result := make([]vec2d.T, len(points))
	const limit = 85.0511287798066
	for i, point := range points {
		lat := math.Max(-limit, math.Min(limit, point[1]))
		result[i] = vec2d.T{6378137 * point[0] * math.Pi / 180,
			6378137 * math.Log(math.Tan(math.Pi/4+lat*math.Pi/360))}
	}
	return result
}

func (p *GCJ02WebMercatorProj) transformFromWGS84(points []vec2d.T) []vec2d.T {
	return p.transformFromGCJ02(newGCJ02Proj(true).transformFromWGS84(points))
}

func (p *GCJ02WebMercatorProj) TransformRectTo(o Proj, rect vec2d.Rect, withPoints int) vec2d.Rect {
	if p.Eq(o) {
		return rect
	}
	return CalculateBBox(p.TransformTo(o, GenerateEnvelopePoints(rect, withPoints)))
}

func (p *GCJ02WebMercatorProj) Eq(o Proj) bool                       { _, ok := o.(*GCJ02WebMercatorProj); return ok }
func (p *GCJ02WebMercatorProj) AlignBBox(rect vec2d.Rect) vec2d.Rect { return rect }
func (p *GCJ02WebMercatorProj) GetSrsCode() string                   { return "GCJ02:3857" }
func (p *GCJ02WebMercatorProj) GetDef() string                       { return p.GetSrsCode() }
func (p *GCJ02WebMercatorProj) ToString() string                     { return p.GetSrsCode() }
func (p *GCJ02WebMercatorProj) IsLatLong() bool                      { return false }
func (p *GCJ02WebMercatorProj) IsAxisOrderNE() bool                  { return false }
