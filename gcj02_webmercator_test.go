package geo

import (
	"math"
	"testing"

	vec2d "github.com/flywave/go3d/float64/vec2"
)

func TestGCJ02WebMercatorLandmarks(t *testing.T) {
	wgs := NewProj("EPSG:4326")
	web := NewProj("GCJ02:3857")
	// GCJ02 fixture coordinates computed with gcoord 1.0.7; spherical Mercator
	// fixture meters computed independently with R*asinh(tan(latitude)).
	for _, tc := range []struct {
		name        string
		wgs, meters vec2d.T
	}{
		{"beijing", vec2d.T{116.46706996, 39.99188446}, vec2d.T{12965738.109612077, 4864954.606322535}},
		{"shanghai", vec2d.T{121.4947167, 31.2416694}, vec2d.T{13525225.509552354, 3663914.331237133}},
		{"guangzhou", vec2d.T{113.31889, 23.10889}, vec2d.T{12615203.357322499, 2644876.298279130}},
		{"wuhan", vec2d.T{114.29694, 30.54694}, vec2d.T{12724081.617749535, 3573733.995621389}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := wgs.TransformTo(web, []vec2d.T{tc.wgs})[0]
			if math.Hypot(got[0]-tc.meters[0], got[1]-tc.meters[1]) > 0.05 {
				t.Fatalf("got %v, want %v", got, tc.meters)
			}
			back := web.TransformTo(wgs, []vec2d.T{got})[0]
			if math.Abs(back[0]-tc.wgs[0]) > 1e-6 || math.Abs(back[1]-tc.wgs[1]) > 1e-6 {
				t.Fatalf("inverse: %v", back)
			}
			merc := NewProj("EPSG:3857")
			via := merc.TransformTo(web, wgs.TransformTo(merc, []vec2d.T{tc.wgs}))[0]
			if math.Hypot(via[0]-got[0], via[1]-got[1]) > 0.001 {
				t.Fatalf("3857 path mismatch: %v", via)
			}
		})
	}
	if web.Eq(NewProj("EPSG:GCJ02MC")) {
		t.Fatal("new and legacy projections must differ")
	}
	if GetLatLongProj(web).GetSrsCode() != "EPSG:GCJ02" {
		t.Fatal("incorrect geographic projection")
	}
}
func TestGCJ02WebMercatorGeographicAndPoles(t *testing.T) {
	web := NewProj("GCJ02:3857")
	gcj := NewProj("EPSG:GCJ02")
	points := []vec2d.T{{121.49916796394122, 31.23967156200282}, {0, 90}, {0, -90}}
	meters := gcj.TransformTo(web, points)
	for _, point := range meters {
		if math.IsNaN(point[1]) || math.IsInf(point[1], 0) || math.Abs(point[1]) > 20037508.343 {
			t.Fatalf("invalid Mercator: %v", point)
		}
	}
	back := web.TransformTo(gcj, meters)[0]
	if math.Abs(back[0]-points[0][0]) > 1e-9 || math.Abs(back[1]-points[0][1]) > 1e-9 {
		t.Fatalf("GCJ02 roundtrip %v", back)
	}
}
