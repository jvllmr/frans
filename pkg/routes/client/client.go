package clientRoutes

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jvllmr/frans/pkg/config"
	"github.com/jvllmr/frans/pkg/ent/ticket"
	routesUtil "github.com/jvllmr/frans/pkg/routes"
	"github.com/tidwall/gjson"
)

//go:embed assets/*
var clientFiles embed.FS

//go:embed index.html.tmpl
var indexFileContent string

func loadManifest() gjson.Result {
	var jsonData, _ = clientFiles.ReadFile("assets/.vite/manifest.json")
	return gjson.ParseBytes(jsonData)
}

func setIndexTemplate(r *gin.Engine, configValue config.Config) *gin.Engine {
	var manifest gjson.Result
	if !configValue.DevMode {
		manifest = loadManifest()
	} else {
		manifest = gjson.Parse("{}")
	}

	assetUrl := func(filePath string) string {

		if configValue.DevMode {
			return fmt.Sprintf("http://localhost:3000/%s", filePath)
		}

		filePathResult := manifest.Get(fmt.Sprintf("%s.file", gjson.Escape(filePath)))

		return fmt.Sprintf("%s/static/%s", configValue.RootPath, filePathResult.String())
	}
	asset := func(filePath string) string {
		return assetUrl((filePath))
	}

	css := func(filePath string) []string {
		var urls []string
		css_files_result := manifest.Get(fmt.Sprintf("%s.css", gjson.Escape(filePath)))
		if !css_files_result.Exists() || !css_files_result.IsArray() {
			return make([]string, 0)
		}
		for _, css_file_result := range css_files_result.Array() {
			css_file := css_file_result.String()
			if configValue.DevMode {
				urls = append(urls, fmt.Sprintf("http://localhost:3000/%s", css_file))
			} else {
				urls = append(urls, fmt.Sprintf("%s/static/%s", configValue.RootPath, css_file))
			}
		}
		return urls
	}

	imports := func(filePath string) []string {
		var urls []string
		files_result := manifest.Get(fmt.Sprintf("%s.imports", gjson.Escape(filePath)))
		if !files_result.Exists() || !files_result.IsArray() {
			return make([]string, 0)
		}
		for _, file_result := range files_result.Array() {
			file_path := file_result.String()
			urls = append(urls, assetUrl(file_path))

		}
		return urls
	}

	importCSS := func(filePath string) []string {
		var urls []string
		files_result := manifest.Get(fmt.Sprintf("%s.imports", gjson.Escape(filePath)))
		if !files_result.Exists() || !files_result.IsArray() {
			return make([]string, 0)
		}
		for _, file_result := range files_result.Array() {
			file_path := file_result.String()
			urls = append(urls, css(file_path)...)
		}
		return urls
	}

	tmpl := template.Must(template.New("index").Funcs(template.FuncMap{
		"assetUrl":   assetUrl,
		"asset":      asset,
		"css":        css,
		"imports":    imports,
		"importsCSS": importCSS,
	}).Parse(indexFileContent))
	r.SetHTMLTemplate(tmpl)

	return r
}

func SetupClientRoutes(r *gin.Engine, rGroup *gin.RouterGroup, configValue config.Config) {
	rGroup.GET("/s/:id", func(c *gin.Context) {
		id := c.Param("id")
		uuid, err := uuid.Parse(id)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}
		targetTicket, _ := config.DBClient.Ticket.Query().
			Where(ticket.ID(uuid)).
			First(c.Request.Context())
		if targetTicket != nil {
			c.Redirect(
				http.StatusPermanentRedirect,
				fmt.Sprintf("/share/ticket/%s", targetTicket.ID.String()),
			)
			return
		}

	})

	setIndexTemplate(r, configValue)
	staticFiles, _ := fs.Sub(clientFiles, "assets")
	rGroup.StaticFS("/static", http.FS(staticFiles))

	r.NoRoute(routesUtil.AuthMiddleware(configValue, true), func(c *gin.Context) {
		// Fallback to index.html for React Router
		c.HTML(http.StatusOK, "index", gin.H{
			"rootPath":                           configValue.RootPath,
			"devMode":                            configValue.DevMode,
			"maxFiles":                           configValue.MaxFiles,
			"maxSizes":                           configValue.MaxSizes,
			"defaultExpiryTotalDays":             configValue.DefaultExpiryTotalDays,
			"defaultExpiryTotalDownloads":        configValue.DefaultExpiryTotalDownloads,
			"defaultExpiryDaysSinceLastDownload": configValue.DefaultExpiryDaysSinceLastDownload,
		})
	})

}
