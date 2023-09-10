#Templ
A tool to render templates, and download and update repositories of templates.

![logo](.github/images/templ-logo-smaller.png)

# Environment variables
TEMPL_DIR - Defaults to ~/.config/templ. Directory that stores template git repositories. You run lists against this dir.

# Uses
`templ -fetch https://github.com/PlayTechniuque/templ-templates.git` - download a github repository to your templates directory

`templ -list` - list all available templates

`templ templatename` - display the contents of a template file to stdout.

`templ templatename=variablesfile.yaml` - hydrate a template file using variablesfile.yaml. 

```variablesfile.yaml
---
Variable: value
OtherVariable: othervalue
```

`templ templatename | templ KEY=VALUE` - pipeline a render operation! You do not need a config file, you can pipe through templ itself.
