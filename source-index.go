package main

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2017 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"text/template"

	"pkg.re/essentialkaos/ek.v7/arg"
	"pkg.re/essentialkaos/ek.v7/env"
	"pkg.re/essentialkaos/ek.v7/fmtc"
	"pkg.re/essentialkaos/ek.v7/fsutil"
	"pkg.re/essentialkaos/ek.v7/sortutil"
	"pkg.re/essentialkaos/ek.v7/timeutil"
	"pkg.re/essentialkaos/ek.v7/usage"
	"pkg.re/essentialkaos/ek.v7/usage/update"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	APP  = "SourceIndex"
	VER  = "0.0.1"
	DESC = "Utility for generating index for source archives"
)

const (
	ARG_OUTPUT   = "o:output"
	ARG_TEMPLATE = "t:template"
	ARG_NO_COLOR = "nc:no-color"
	ARG_HELP     = "h:help"
	ARG_VER      = "v:version"
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

var argMap = arg.Map{
	ARG_OUTPUT:   {Value: "index.html"},
	ARG_TEMPLATE: {Value: "default.tpl"},
	ARG_NO_COLOR: {Type: arg.BOOL},
	ARG_HELP:     {Type: arg.BOOL, Alias: "u:usage"},
	ARG_VER:      {Type: arg.BOOL, Alias: "ver"},
}

// ////////////////////////////////////////////////////////////////////////////////// //

func main() {
	args, errs := arg.Parse(argMap)

	if len(errs) != 0 {
		for _, err := range errs {
			printError(err.Error())
		}

		os.Exit(1)
	}

	if arg.GetB(ARG_NO_COLOR) {
		fmtc.DisableColors = true
	}

	if arg.GetB(ARG_VER) {
		showAbout()
		return
	}

	if arg.GetB(ARG_HELP) || len(args) == 0 {
		showUsage()
		return
	}

	process(args[0])
}

// process start processing
func process(dir string) {
	err := checkDir(dir)

	if err != nil {
		printError(err.Error())
		os.Exit(1)
	}

	index := buildIndex(dir)

	err = export(index)

	if err != nil {
		printError(err.Error())
		os.Exit(2)
	}

	projects, releases := index.Stats()

	fmtc.Printf(
		"{g}Index for %d projects and %d releases successfully generated as {g*}%s{!}\n",
		projects, releases, arg.GetS(ARG_OUTPUT),
	)
}

// checkDir check directory
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

// buildIndex build index with info about all projects in directory
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

// getReleases read given directory and return slice with info about releases
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

// parseSourceName parse source name and return version and source info
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

// releaseMapToSlice convert map with releases to sorted slice
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

// export render template with inforamtion from index and save as file
func export(index *Index) error {
	templateFile := getTemplateFile()
	outputFile := arg.GetS(ARG_OUTPUT)

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

// getTemplateFile return path to template file
func getTemplateFile() string {
	template := arg.GetS(ARG_TEMPLATE)

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
	fmtc.Printf("{r}"+f+"{!}\n", a...)
}

// printWarn prints warning message to console
func printWarn(f string, a ...interface{}) {
	fmtc.Printf("{y}"+f+"{!}\n", a...)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Stats return number of projects and releases in index
func (i *Index) Stats() (int, int) {
	var releases int

	for _, project := range i.Projects {
		releases += len(project.Releases)
	}

	return len(i.Projects), releases
}

// ////////////////////////////////////////////////////////////////////////////////// //

func showUsage() {
	info := usage.NewInfo("", "dir")

	info.AddOption(ARG_OUTPUT, "Output file {s-}(index.html by default){!}", "file")
	info.AddOption(ARG_TEMPLATE, "Template {s-}(template.tpl by default){!}", "file")
	info.AddOption(ARG_NO_COLOR, "Disable colors in output")
	info.AddOption(ARG_HELP, "Show this help message")
	info.AddOption(ARG_VER, "Show version")

	info.Render()
}

func showAbout() {
	about := &usage.About{
		App:           APP,
		Version:       VER,
		Desc:          DESC,
		Year:          2006,
		Owner:         "ESSENTIAL KAOS",
		License:       "Essential Kaos Open Source License <https://essentialkaos.com/ekol>",
		UpdateChecker: usage.UpdateChecker{"essentialkaos/source-index", update.GitHubChecker},
	}

	about.Render()
}
