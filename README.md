# Templ
A tool to render go text templates. The sweetness is that you can store your text templates in a git repository,
and use templ to download and use that repository of templates.

Update templates with your text editor, commit to git with your usual workflow, then `templ -u` to download the updates.

Want to check that your templates are written and parsed correctly before committing them to git? Awesome, templ can be 
piped into:

cat foo.tmpl | templ VAR=VALUE

![logo](.github/images/templ-logo-smaller.png)

# Environment variables
TEMPL_DIR - Defaults to ~/.config/templ. Directory that stores template git repositories. You run lists against this dir.

# Uses
## Retrieve Templates Stored in a Git Repository
`templ -f https://github.com/PlayTechniuque/templ-templates.git` - download a github repository to your templates directory.
You can set the templates directory with an environment variable TEMPL_DIR or use the default (currently ~/.config/templ)

## List your templates
`templ -l` - list all downloaded templates
`templ templatename` - display the contents of a template file to stdout. It's like cat, but it has partial file matching,
so if you have a template file in a directory structure `foo/bar/bam.yaml` then you can use `templ foo` or `templ bam.yaml`
or `templ bar/bam`. 

If you have two files `foo/bar/bam.yaml` and `zee/zye/bam.yaml`, then `templ bam` will show both of the `bam.yaml` files.

## Rendering templates
You have two options for rendering templates. The first and simplest is to put the template file on
stdout and then pipe that template through templ itself, replacing variables with values:
`templ templatename | templ KEY=VALUE` - pipeline a render operation! You do not need a config file, you can pipe through
templ itself. This is best suited for replacing simple variables, not using template logic

`templ templatename=variablesfile.yaml` - hydrate a template file using variablesfile.yaml. This gives you more 
flexibility in composing more complex templates.

```variablesfile.yaml
---
Variable: value
OtherVariable: othervalue
```
