package container

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/bytesparadise/libasciidoc"
	asciidocconf "github.com/bytesparadise/libasciidoc/pkg/configuration"
	"github.com/shemanaev/updtr/internal/config"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

func (c *Client) getChangelog(typ config.ChangelogType, url string) (string, error) {
	switch typ {
	case config.Plaintext:
		res, err := http.Get(url)
		if err != nil {
			return "", err
		}

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(`<pre>%s</pre>`, resBody), nil

	case config.Html:
		res, err := http.Get(url)
		if err != nil {
			return "", err
		}

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return "", err
		}
		return string(resBody), nil

	case config.Markdown:
		res, err := http.Get(url)
		if err != nil {
			return "", err
		}

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return "", err
		}

		md := goldmark.New(
			goldmark.WithExtensions(extension.GFM),
		)

		var buf bytes.Buffer
		if err := md.Convert(resBody, &buf); err != nil {
			return "", err
		}

		return buf.String(), nil

	case config.Asciidoc:
		res, err := http.Get(url)
		if err != nil {
			return "", err
		}

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return "", err
		}

		output := &strings.Builder{}
		content := bytes.NewReader(resBody)
		libasciidoc.Convert(content, output, asciidocconf.NewConfiguration())

		return output.String(), nil

	case config.Github:
		match := reGithub.FindSubmatch([]byte(url))
		if len(match) == 0 {
			return "", fmt.Errorf("not valid github link: %s", url)
		}

		url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", string(match[1]), string(match[2]))
		res, err := http.Get(url)
		if err != nil {
			return "", err
		}

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return "", err
		}

		type ghRelease struct {
			Body string
		}

		var release ghRelease
		if err := json.Unmarshal(resBody, &release); err != nil {
			return "", err
		}

		md := goldmark.New(
			goldmark.WithExtensions(extension.GFM),
		)

		var buf bytes.Buffer
		if err := md.Convert([]byte(release.Body), &buf); err != nil {
			return "", err
		}

		return buf.String(), nil
	}

	return "", nil
}
