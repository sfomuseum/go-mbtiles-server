module github.com/sfomuseum/go-mbtiles-server

go 1.16

require (
	github.com/aaronland/go-http-server v0.0.7
	github.com/aaronland/go-mimetypes v0.0.1
	github.com/tilezen/go-tilepacks v0.0.0-20210602151652-0147f7fb6fd7
)

replace github.com/tilezen/go-tilepacks v0.0.0-20210602151652-0147f7fb6fd7 => github.com/sfomuseum/go-tilepacks v0.0.0-20210818170442-a65e38f5c453
