module auroraTags/mmo_game

go 1.19

require (
	aurora-v1.0 v0.0.0-00010101000000-000000000000
	google.golang.org/protobuf v1.30.0
)

// 依赖工具包替换为本地 注意替换的包下面必须有 go.mod 文件
replace aurora-v1.0 => ../aurora-1.0
