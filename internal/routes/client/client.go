package clientRoutes

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jvllmr/frans/internal/config"
	"github.com/jvllmr/frans/internal/ent"
	"github.com/jvllmr/frans/internal/ent/grant"
	"github.com/jvllmr/frans/internal/ent/ticket"
	"github.com/jvllmr/frans/internal/middleware"
	"github.com/jvllmr/frans/internal/oidc"

	"github.com/tidwall/gjson"
)

//go:embed all:assets/*
var clientFiles embed.FS

//go:embed index.gohtml
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

type clientController struct {
	config config.Config
	db     *ent.Client
}

func (cc *clientController) redirectShareLink(c *gin.Context) {
	id := c.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	targetTicket := cc.db.Ticket.Query().
		Where(ticket.ID(uuid)).
		FirstX(c.Request.Context())
	if targetTicket != nil {
		c.Redirect(
			http.StatusPermanentRedirect,
			fmt.Sprintf("%s/share/ticket/%s", cc.config.RootPath, targetTicket.ID.String()),
		)
		return
	}

	targetGrant := cc.db.Grant.Query().
		Where(grant.ID(uuid)).
		FirstX(c.Request.Context())
	if targetGrant != nil {
		c.Redirect(
			http.StatusPermanentRedirect,
			fmt.Sprintf("%s/share/grant/%s", cc.config.RootPath, targetGrant.ID.String()),
		)
		return
	}
}

func SetupClientRoutes(
	r *gin.Engine,
	rGroup *gin.RouterGroup,
	configValue config.Config,
	db *ent.Client,
	oidcProvider *oidc.FransOidcProvider,
) {
	controller := clientController{
		config: configValue,
		db:     db,
	}

	rGroup.GET("/s/:id", controller.redirectShareLink)

	setIndexTemplate(r, configValue)
	staticFiles, _ := fs.Sub(clientFiles, "assets")
	rGroup.StaticFS("/static", http.FS(staticFiles))

	customColorJsonBytes, err := json.Marshal(configValue.CustomColor)
	if err != nil {
		log.Fatalf("Could not generate json from custom color setting.")
	}
	customColorJson := string(customColorJsonBytes)

	r.NoRoute(middleware.Auth(oidcProvider, true), func(c *gin.Context) {
		// Fallback to index.html for React Router
		c.HTML(http.StatusOK, "index", gin.H{
			"rootPath":                              configValue.RootPath,
			"devMode":                               configValue.DevMode,
			"maxFiles":                              configValue.MaxFiles,
			"maxSizes":                              configValue.MaxSizes,
			"defaultExpiryTotalDays":                configValue.DefaultExpiryTotalDays,
			"defaultExpiryTotalDownloads":           configValue.DefaultExpiryTotalDownloads,
			"defaultExpiryDaysSinceLastDownload":    configValue.DefaultExpiryDaysSinceLastDownload,
			"grantDefaultExpiryTotalDays":           configValue.GrantDefaultExpiryTotalDays,
			"grantDefaultExpiryTotalUploads":        configValue.GrantDefaultExpiryTotalUploads,
			"grantDefaultExpiryDaysSinceLastUpload": configValue.GrantDefaultExpiryDaysSinceLastUpload,
			"fransVersion":                          config.FransVersion,
			"color":                                 configValue.Color,
			"customColor":                           customColorJson,
		})
	})

}
