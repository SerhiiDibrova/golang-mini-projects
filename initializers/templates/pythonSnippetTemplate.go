package initializers

import string

class PythonSnippetTemplate:
    def __init__(self, params):
        self.params = params
        self.template_string = string.Template('print("Hello, $Name!")')

    def render(self):
        return self.template_string.substitute(self.params)

def python_snippet_template(params):
    template = PythonSnippetTemplate(params)
    return template.render()