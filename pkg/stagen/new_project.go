package stagen

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/pixality-inc/golang-core/http"
	"github.com/pixality-inc/golang-core/logger"
	"gopkg.in/yaml.v3"
)

type newProjectConfig struct {
	Logger logger.YamlConfig `yaml:"logger"`
	Stagen ConfigYaml        `yaml:"stagen"`
	Site   SiteConfigYaml    `yaml:"site"`
}

//go:embed assets/stagen_64.png
var logoPng []byte

func (s *Impl) NewProject(ctx context.Context, name string) error {
	log := s.log.GetLogger(ctx)

	log.Infof("Creating new project: %s", name)

	hasGit := s.git.HasGit(ctx)

	sourceDir := s.workDir()

	databasesDir := s.databasesDir()
	extDir := s.extensionsDir()
	pagesDir := s.pagesDir()
	publicDir := s.publicDir()
	publicCssDir := filepath.Join(publicDir, "css")
	publicJsDir := filepath.Join(publicDir, "js")
	templatesDir := s.templatesDir()
	themesDir := s.themesDir()

	dirs := []string{
		sourceDir,
		databasesDir,
		extDir,
		pagesDir,
		publicDir,
		publicCssDir,
		publicJsDir,
		templatesDir,
		themesDir,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	cfg := &newProjectConfig{
		Logger: logger.YamlConfig{
			LevelValue:            "info",
			FormatValue:           "text",
			TimestampValue:        true,
			ColorsValue:           true,
			StacktraceValue:       false,
			StacktraceErrorsValue: false,
		},
		Stagen: ConfigYaml{
			EnvValue: "dev",
			HttpValue: http.ConfigYaml{
				HostValue:            "127.0.0.1",
				PortValue:            8001,
				ShutdownTimeoutValue: 10 * time.Second,
			},
			SettingsValue: ConfigSettingsYaml{
				UseUriHtmlFileExtensionValue: false,
			},
			DirsValue: ConfigDirsYaml{
				WorkValue:       "",
				BuildValue:      "",
				DatabasesValue:  "",
				ExtensionsValue: "",
				ThemesValue:     "",
				TemplatesValue:  "",
				PagesValue:      "",
				PublicValue:     "",
			},
		},
		Site: SiteConfigYaml{
			BaseUrlValue:     "http://127.0.0.1:8001",
			NameValue:        name,
			DescriptionValue: "",
			LangValue:        "en",
			AuthorValue: SiteConfigAuthorYaml{
				NameValue:    "",
				EmailValue:   "",
				WebsiteValue: "",
			},
			LogoValue: SiteConfigLogoYaml{
				UrlValue: "/logo.png",
			},
			CopyrightValue: SiteConfigCopyrightYaml{
				YearValue:   time.Now().Year(),
				TitleValue:  "Stagen",
				RightsValue: "All rights reserved.",
			},
			ExtensionsValue: nil,
			AggDictsValue:   nil,
			GeneratorsValue: nil,
			TemplateValue: SiteConfigTemplateYaml{
				ThemeValue:         "default",
				DefaultLayoutValue: "",
				VariablesValue:     nil,
				ImportsValue:       nil,
				IncludesValue:      nil,
				ExtrasValue:        nil,
			},
		},
	}

	configYaml, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	styleCss := []byte(``)

	appJs := []byte(``)

	htaccess := []byte(`# mod_rewrite
RewriteEngine On

# check for file.html
RewriteCond %{REQUEST_FILENAME} !-f
RewriteCond %{REQUEST_FILENAME}.html -f
RewriteRule ^(.+)$ $1.html [L]

# error pages
ErrorDocument 404 /404.html
ErrorDocument 500 /50x.html
ErrorDocument 502 /50x.html
ErrorDocument 503 /50x.html
`)
	robotsTxt := []byte(`User-agent: *
Allow: /
`)

	readmeMd := []byte("# Stagen website\n\n")

	indexMd := fmt.Appendf(nil, `---
title: %s
---

Welcome to my cool stagen website!

![](/logo.png)`, strconv.Quote(name))

	err50xHtml := []byte(`<h1>Internal Server Error</h1>`)

	err404Html := []byte(`<h1>Not Found</h1>`)

	mainMenuYaml := []byte(`---
name: main_menu
data:
  - page: index
    title: Home
`)

	gitignore := []byte(`.DS_Store
.idea

/build`)

	filesToWrite := map[string][]byte{
		filepath.Join(sourceDir, "config.yaml"): configYaml,
		filepath.Join(sourceDir, "README.md"):   readmeMd,
		filepath.Join(sourceDir, ".gitignore"):  gitignore,

		filepath.Join(publicDir, "logo.png"):   logoPng,
		filepath.Join(publicDir, ".htaccess"):  htaccess,
		filepath.Join(publicDir, "robots.txt"): robotsTxt,

		filepath.Join(publicCssDir, "style.css"): styleCss,

		filepath.Join(publicJsDir, "app.js"): appJs,

		filepath.Join(pagesDir, "index.md"): indexMd,
		filepath.Join(pagesDir, "404.html"): err404Html,
		filepath.Join(pagesDir, "50x.html"): err50xHtml,

		filepath.Join(databasesDir, "main_menu.yaml"): mainMenuYaml,
	}

	for filename, content := range filesToWrite {
		//nolint:gosec
		if err = os.WriteFile(filename, content, os.ModePerm); err != nil {
			return fmt.Errorf("failed to write file %s: %w", filename, err)
		}
	}

	if hasGit {
		if err = s.git.Init(ctx, sourceDir); err != nil {
			return fmt.Errorf("failed to init git: %w", err)
		}

		if err = s.git.SubmoduleAdd(ctx, sourceDir, "https://github.com/Stagens/theme-default.git", "themes/default"); err != nil {
			return fmt.Errorf("failed to add themes/default git submodule: %w", err)
		}
	} else {
		log.Errorf("Git is not installed")
	}

	return nil
}
