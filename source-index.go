package main

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2022 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/essentialkaos/ek/v12/env"
	"github.com/essentialkaos/ek/v12/fmtc"
	"github.com/essentialkaos/ek/v12/fsutil"
	"github.com/essentialkaos/ek/v12/options"
	"github.com/essentialkaos/ek/v12/sortutil"
	"github.com/essentialkaos/ek/v12/timeutil"
	"github.com/essentialkaos/ek/v12/usage"
	"github.com/essentialkaos/ek/v12/usage/completion/bash"
	"github.com/essentialkaos/ek/v12/usage/completion/fish"
	"github.com/essentialkaos/ek/v12/usage/completion/zsh"
	"github.com/essentialkaos/ek/v12/usage/update"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	APP  = "SourceIndex"
	VER  = "0.3.3"
	DESC = "Utility for generating index for source archives"
)

const (
	OPT_OUTPUT   = "o:output"
	OPT_TEMPLATE = "t:template"
	OPT_NO_COLOR = "nc:no-color"
	OPT_HELP     = "h:help"
	OPT_VER      = "v:version"

	OPT_COMPLETION = "completion"
)

// ////////////////////////////////////////////////////////////////////////////////// //

type Index struct {
	Projects []*Project
}

type Project struct {
	Name     string
	Releases []*Release
}

type Release struct {
	Version string
	Sources []*Source
	Date    string
	Latest  bool
}

type Source struct {
	File string
	Ext  string
}

// ////////////////////////////////////////////////////////////////////////////////// //

type ReleaseSlice []*Release

func (s ReleaseSlice) Len() int      { return len(s) }
func (s ReleaseSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ReleaseSlice) Less(i, j int) bool {
	return sortutil.VersionCompare(s[i].Version, s[j].Version)
}

type ProjectSlice []*Project

func (s ProjectSlice) Len() int      { return len(s) }
func (s ProjectSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ProjectSlice) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

type SourceSlice []*Source

func (s SourceSlice) Len() int      { return len(s) }
func (s SourceSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s SourceSlice) Less(i, j int) bool {
	return s[i].Ext < s[j].Ext
}

// ////////////////////////////////////////////////////////////////////////////////// //

var optMap = options.Map{
	OPT_OUTPUT:   {Value: "index.html"},
	OPT_TEMPLATE: {Value: "default.tpl"},
	OPT_NO_COLOR: {Type: options.BOOL},
	OPT_HELP:     {Type: options.BOOL, Alias: "u:usage"},
	OPT_VER:      {Type: options.BOOL, Alias: "ver"},

	OPT_COMPLETION: {},
}

// ////////////////////////////////////////////////////////////////////////////////// //

func main() {
	args, errs := options.Parse(optMap)

	if len(errs) != 0 {
		for _, err := range errs {
			printError(err.Error())
		}

		os.Exit(1)
	}

	if options.Has(OPT_COMPLETION) {
		genCompletion()
	}

	if options.GetB(OPT_NO_COLOR) {
		fmtc.DisableColors = true
	}

	if options.GetB(OPT_VER) {
		showAbout()
		return
	}

	if options.GetB(OPT_HELP) || len(args) == 0 {
		showUsage()
		return
	}

	process(args.Get(0).Clean().String())
}

// process starts processing
func process(dir string) {
	err := checkDir(dir)

	if err != nil {
		printErrorAndExit(err.Error())
	}

	index := buildIndex(dir)

	err = export(index)

	if err != nil {
		printErrorAndExit(err.Error())
	}

	projects, releases := index.Stats()

	fmtc.Printf(
		"{g}Index for %d projects and %d releases successfully generated as {g*}%s{!}\n",
		projects, releases, options.GetS(OPT_OUTPUT),
	)
}

// checkDir checks directory
func checkDir(dir string) error {
	if !fsutil.IsExist(dir) {
		return fmt.Errorf("Directory %s doesn't exist", dir)
	}

	if !fsutil.IsReadable(dir) {
		return fmt.Errorf("Directory %s is not readable", dir)
	}

	if !fsutil.IsExecutable(dir) {
		return fmt.Errorf("Directory %s is not executable", dir)
	}

	if fsutil.IsEmptyDir(dir) {
		return fmt.Errorf("Directory %s is empty", dir)
	}

	return nil
}

// buildIndex builds index with info about all projects in directory
func buildIndex(dir string) *Index {
	var index = &Index{}

	projects := fsutil.List(dir, true, fsutil.ListingFilter{Perms: "DRX"})

	if len(projects) == 0 {
		return index
	}

	for _, projectName := range projects {
		project := &Project{
			Name:     projectName,
			Releases: getReleases(projectName, dir+"/"+projectName),
		}

		if len(project.Releases) == 0 {
			continue
		}

		index.Projects = append(index.Projects, project)
	}

	sort.Sort(ProjectSlice(index.Projects))

	return index
}

// getReleases reads given directory and return slice with info about releases
func getReleases(project, dir string) []*Release {
	var releases map[string]*Release

	sources := fsutil.List(dir, true, fsutil.ListingFilter{Perms: "FR"})

	if len(sources) == 0 {
		return []*Release{}
	}

	releases = make(map[string]*Release)

	for _, sourceName := range sources {
		version, source := parseSourceName(project, sourceName)

		if version == "current" || version == "" {
			continue
		}

		release, ok := releases[version]

		if !ok {
			release = &Release{Version: version, Sources: []*Source{}}
			releases[version] = release
		}

		if release.Date == "" {
			cd, _ := fsutil.GetMTime(dir + "/" + sourceName)
			release.Date = timeutil.Format(cd, "%Y/%m/%d")
		}

		release.Sources = append(release.Sources, source)
	}

	if len(releases) == 0 {
		return []*Release{}
	}

	return releaseMapToSlice(releases)
}

// parseSourceName parses source name and return version and source info
func parseSourceName(project, name string) (string, *Source) {
	verIndex := strings.LastIndex(name, "-")

	if verIndex == -1 {
		return "", nil
	}

	verAndExt := name[verIndex+1:]

	var (
		version string
		ext     string
	)

	switch {
	case strings.HasSuffix(verAndExt, ".zip"):
		version = strings.Replace(verAndExt, ".zip", "", -1)
		ext = "ZIP"

	case strings.HasSuffix(verAndExt, ".7z"):
		version = strings.Replace(verAndExt, ".7z", "", -1)
		ext = "7Z"

	case strings.HasSuffix(verAndExt, ".tar.bz2"):
		version = strings.Replace(verAndExt, ".tar.bz2", "", -1)
		ext = "TAR.BZ2"

	case strings.HasSuffix(verAndExt, ".tbz2"):
		version = strings.Replace(verAndExt, ".tbz2", "", -1)
		ext = "TAR.BZ2"

	case strings.HasSuffix(verAndExt, ".tar.gz"):
		version = strings.Replace(verAndExt, ".tar.gz", "", -1)
		ext = "TAR.GZ"

	case strings.HasSuffix(verAndExt, ".tgz"):
		version = strings.Replace(verAndExt, ".tgz", "", -1)
		ext = "TAR.GZ"

	case strings.HasSuffix(verAndExt, ".tar.xz"):
		version = strings.Replace(verAndExt, ".tar.xz", "", -1)
		ext = "TAR.XZ"

	case strings.HasSuffix(verAndExt, ".txz"):
		version = strings.Replace(verAndExt, ".txz", "", -1)
		ext = "TAR.XZ"
	}

	return version, &Source{File: project + "/" + name, Ext: ext}
}

// releaseMapToSlice converts map with releases to sorted slice
func releaseMapToSlice(releases map[string]*Release) []*Release {
	var result []*Release

	for _, release := range releases {
		sort.Sort(SourceSlice(release.Sources))
		result = append(result, release)
	}

	sort.Sort(sort.Reverse(ReleaseSlice(result)))

	result[0].Latest = true

	return result
}

// export renders template with inforamtion from index and save as file
func export(index *Index) error {
	templateFile := getTemplateFile()
	outputFile := options.GetS(OPT_OUTPUT)

	if templateFile == "" {
		return fmt.Errorf("Can't use given template")
	}

	if fsutil.IsExist(outputFile) {
		err := os.Remove(outputFile)

		if err != nil {
			return err
		}
	}

	fd, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}

	defer fd.Close()

	tpl, err := ioutil.ReadFile(templateFile)

	if err != nil {
		return err
	}

	t := template.New("template")
	t, err = t.Parse(string(tpl[:]))

	return t.Execute(fd, index)
}

// getTemplateFile returns path to template file
func getTemplateFile() string {
	template := options.GetS(OPT_TEMPLATE)

	if fsutil.CheckPerms("FR", template) {
		return template
	}

	gopath := env.Get().GetS("GOPATH")

	template = gopath + "/src/github.com/essentialkaos/source-index/templates/" + template

	if fsutil.CheckPerms("FR", template) {
		return template
	}

	return ""
}

// printError prints error message to console
func printError(f string, a ...interface{}) {
	fmtc.Fprintf(os.Stderr, "{r}"+f+"{!}\n", a...)
}

// printError prints warning message to console
func printWarn(f string, a ...interface{}) {
	fmtc.Fprintf(os.Stderr, "{y}"+f+"{!}\n", a...)
}

// printErrorAndExit prints error mesage and exit with exit code 1
func printErrorAndExit(f string, a ...interface{}) {
	printError(f, a...)
	os.Exit(1)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Stats returns number of projects and releases in index
func (i *Index) Stats() (int, int) {
	var releases int

	for _, project := range i.Projects {
		releases += len(project.Releases)
	}

	return len(i.Projects), releases
}

// ////////////////////////////////////////////////////////////////////////////////// //

// showUsage prints usage info
func showUsage() {
	genUsage().Render()
}

// genUsage
func genUsage() *usage.Info {
	info := usage.NewInfo("", "dir")

	info.AddOption(OPT_OUTPUT, "Output file {s-}(index.html by default){!}", "file")
	info.AddOption(OPT_TEMPLATE, "Template {s-}(template.tpl by default){!}", "file")
	info.AddOption(OPT_NO_COLOR, "Disable colors in output")
	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VER, "Show version")

	return info
}

// genCompletion generates completion for different shells
func genCompletion() {
	info := genUsage()

	switch options.GetS(OPT_COMPLETION) {
	case "bash":
		fmt.Printf(bash.Generate(info, "source-index"))
	case "fish":
		fmt.Printf(fish.Generate(info, "source-index"))
	case "zsh":
		fmt.Printf(zsh.Generate(info, optMap, "source-index"))
	default:
		os.Exit(1)
	}

	os.Exit(0)
}

// showAbout prints basic info about app
func showAbout() {
	about := &usage.About{
		App:           APP,
		Version:       VER,
		Desc:          DESC,
		Year:          2006,
		Owner:         "ESSENTIAL KAOS",
		License:       "Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>",
		UpdateChecker: usage.UpdateChecker{"essentialkaos/source-index", update.GitHubChecker},
	}

	about.Render()
}
