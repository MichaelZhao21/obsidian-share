package src

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// CloneOrPullRepo clones a repository if it doesn't exist, or pulls if it does.
// Requires SSH_PRIVATE_KEY environment variable to be set.
func CloneOrPullRepo(repoURL, localPath string) error {
	// Get SSH key from environment variable
	sshKeyFile := os.Getenv("SSH_PRIVATE_KEY")
	if sshKeyFile == "" {
		return errors.New("SSH_PRIVATE_KEY environment variable is not set")
	}

	// Read the SSH key from a file
	publicKeys, err := ssh.NewPublicKeysFromFile("git", sshKeyFile, "")
	if err != nil {
		return err
	}

	// Make sure the local path exists
	err = os.MkdirAll(localPath, os.ModePerm)
	if err != nil {
		return err
	}

	// Check if repository already exists
	if _, err := os.Stat(filepath.Join(localPath, ".git")); err == nil {
		// Repository exists, perform pull
		repo, err := git.PlainOpen(localPath)
		if err != nil {
			return err
		}

		worktree, err := repo.Worktree()
		if err != nil {
			return err
		}

		err = worktree.Pull(&git.PullOptions{
			RemoteName: "origin",
			Auth:       publicKeys,
		})
		if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
			return err
		}
		return nil
	}

	// Repository doesn't exist, perform clone
	_, err = git.PlainClone(localPath, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout, // Shows clone progress
		Auth:     publicKeys,
	})
	return err
}

// Convert all given md files to HTML
func ConvAllToHtml(files []FileLink, outPath string) error {
	// Create the output directory if it doesn't exist
	err := os.MkdirAll(outPath, os.ModePerm)
	if err != nil {
		return err
	}

	for _, file := range files {
		err := ConvToHtml(file.Slug, filepath.Join("./files/raw", file.Filename), outPath)
		if err != nil {
			return err
		}
	}
	return nil
}

// Convert a single file to HTML
func ConvToHtml(slug string, fileName string, outPath string) error {
	// Get the file info
	info, err := os.Stat(fileName + ".md")
	if err != nil {
		return err
	}

	// Get the HTML file path
	htmlPath := filepath.Join(outPath, filepath.Base(fileName)+".html")

	// Check if HTML file exists and is older than the markdown file
	htmlInfo, err := os.Stat(htmlPath)
	if err == nil && info.ModTime().Before(htmlInfo.ModTime()) {
		return nil
	}
	log.Println("Updating", fileName)

	// Read in the markdown file
	md, err := os.ReadFile(fileName + ".md")
	if err != nil {
		return err
	}

	// Remove all tags from the markdown file (^#text)
	regex := regexp.MustCompile("^#.*")
	mdClean := regex.ReplaceAllString(string(md), "")

	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(mdClean))

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	// render the markdown into HTML
	htmlOutput := markdown.Render(doc, renderer)

	// Load the HTML template from a file
	template, err := os.ReadFile("template.html")
	if err != nil {
		return err
	}

	// Replace the template's body with the HTML output
	replaced := strings.Replace(string(template), "{CONTENT}", string(htmlOutput), -1)
	replaced = strings.Replace(replaced, "{TITLE}", filepath.Base(fileName), -1)
	replaced = strings.Replace(replaced, "{TITLE_DISPLAY}", filepath.Base(fileName), -1)
	replaced = strings.Replace(replaced, "{SLUG}", slug, -1)

	// Write the HTML file
	return os.WriteFile(htmlPath, []byte(replaced), 0644)
}
