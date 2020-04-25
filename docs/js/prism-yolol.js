Prism.languages.yolol = {
    'comment': /\/\/.+/,
    'string': /"[^"]*"/,
    'keyword': /(?i)\b(if|then|else|end|goto)\b/,
    'operator': /(?i)\b(and|or|not)\b/,
    'function': /[a-z0-9_]+(?=\()/i,
    'variable': /:?[a-zA-Z]+[a-zA-Z0-9_]*/,
    'constant': /[0-9]+(\.[0-9]+)?/
};
