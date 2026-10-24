# GCJ-02 Web Mercator（高德 XYZ）

`NewProj("GCJ02:3857")` 表示先将 WGS84 经纬度转换为 GCJ-02，再进行半径 6378137 米的球面 Web Mercator 投影。此名称是自定义标识，不是 EPSG 注册代码。

高德标准 XYZ 瓦片采用这一组合。已有 `EPSG:GCJ02MC` 使用另一套多项式投影，已有 `gcj02/china-cities.json` 测试数据也依赖该行为；因此新增投影，不修改旧投影及其测试数据。

```go
wgs := geo.NewProj("EPSG:4326")
amap := geo.NewProj("GCJ02:3857")
meters := wgs.TransformTo(amap, []vec2d.T{{116.46706996, 39.99188446}})
original := amap.TransformTo(wgs, meters)
```

转换也支持标准 `EPSG:3857` 与此自定义投影之间的路径。高纬度输入限制在 Web Mercator 的有效范围内。GCJ-02 的适用区域和精度沿用现有转换实现；不能以往返误差替代独立的测绘精度验证。

新增测试包含北京、上海、广州、武汉的固定参考值、标准 3857 路径、逆转换、极区处理，以及旧投影区分。旧坐标转换测试保留不变。
