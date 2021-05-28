package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gomarkdown/markdown/ast"
	mdast "github.com/gomarkdown/markdown/ast"
	mdhtml "github.com/gomarkdown/markdown/html"
	mdparser "github.com/gomarkdown/markdown/parser"

	"github.com/gomarkdown/markdown"
	"github.com/pkg/errors"
)

var docsPath = flag.String("docsPath", "", "Path to docs markdown")
var outputDirPath = flag.String("outputDirPath", "", "Path to output directory")

func main() {
	flag.Parse()
	if err := buildSite(*docsPath, *outputDirPath); err != nil {
		log.Fatal(err)
	}
}

func buildSite(docsPath, outputDirPath string) error {
	os.MkdirAll(outputDirPath, 0744)
	os.MkdirAll(filepath.Join(outputDirPath, "docs"), 0744)

	tmpl, err := template.New("template.html").ParseFiles("template.html")
	if err != nil {
		return errors.Wrapf(err, "template.ParseFiles")
	}

	fileInfos, err := os.ReadDir(docsPath)
	if err != nil {
		return errors.Wrapf(err, "os.ReadDir")
	}

	for _, fi := range fileInfos {
		srcPath := filepath.Join(docsPath, fi.Name())
		dstPath := filepath.Join(outputDirPath, "docs", markdownExtToHtml(fi.Name()))
		if err := markdownToHtml(tmpl, srcPath, dstPath); err != nil {
			return errors.Wrapf(err, "markdownToHtml")
		}
	}

	return nil
}

func copyFile(srcPath, dstPath string) error {
	log.Printf("Copying %s to %s\n", srcPath, dstPath)

	srcFile, err := os.Open(srcPath)
	if err != nil {
		return errors.Wrapf(err, "os.Open")
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return errors.Wrapf(err, "os.Create")
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return errors.Wrapf(err, "io.Copy")
	}

	return nil
}

func markdownToHtml(tmpl *template.Template, srcPath, dstPath string) error {
	log.Printf("Converting markdown '%s' to html '%s'\n", srcPath, dstPath)
	md, err := os.ReadFile(srcPath)
	if err != nil {
		return errors.Wrapf(err, "os.ReadFile")
	}
	extensions := mdparser.CommonExtensions | mdparser.Footnotes | mdparser.AutoHeadingIDs
	parser := mdparser.NewWithExtensions(extensions)
	astNode := markdown.Parse(md, parser)
	rewriteLinks(astNode)

	f, err := os.Create(dstPath)
	if err != nil {
		return errors.Wrapf(err, "os.Create")
	}
	defer f.Close()

	return tmpl.Execute(f, map[string]interface{}{
		"Title":   fmt.Sprintf(strings.ToLower(firstHeading(astNode))),
		"Content": renderMarkdownAsHtml(astNode),
	})
}

func renderMarkdownAsHtml(node mdast.Node) string {
	var options mdhtml.RendererOptions
	renderer := mdhtml.NewRenderer(options)
	content := markdown.Render(node, renderer)
	return string(content)
}

func firstHeading(node mdast.Node) string {
	var heading string
	visitor := func(node mdast.Node, entering bool) mdast.WalkStatus {
		parent := node.GetParent()
		if parent == nil {
			return mdast.GoToNext
		}

		_, ok := parent.(*mdast.Heading)
		if ok {
			heading = string(node.AsLeaf().Literal)
			return mdast.Terminate
		}
		return mdast.GoToNext
	}
	mdast.Walk(node, mdast.NodeVisitorFunc(visitor))
	return heading
}

func rewriteLinks(node ast.Node) {
	visitor := func(node mdast.Node, entering bool) mdast.WalkStatus {
		link, ok := node.(*mdast.Link)
		if !ok {
			return mdast.GoToNext
		}

		if entering {
			u, err := url.Parse(string(link.Destination))
			if err != nil {
				log.Printf("Could not parse URL '%s': %s\n", u, err)
				return mdast.GoToNext
			}

			if !u.IsAbs() {
				rewriteDst := markdownExtToHtml(string(link.Destination))
				link.Destination = []byte(rewriteDst)
			}
		}
		return mdast.GoToNext
	}
	mdast.Walk(node, mdast.NodeVisitorFunc(visitor))
}

func markdownExtToHtml(path string) string {
	ext := filepath.Ext(path)
	if ext != ".md" {
		return path
	}
	return path[0:len(path)-len(ext)] + ".html"
}
