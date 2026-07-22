const goSnippet = `package main

import "fmt"

func hello(who string) {
    fmt.Printf("Hello, %s!\\n", who)
}`;

const pythonSnippet = `def hello(who="world"):
    print(f"Hello, {who}!")`;

const javascriptSnippet = `function hello(who = "world") {
    console.log(\`Hello, \${who}!\`)
}`;

export default function getSnippet(name: string): string {

    const snippets: Record<string, string> = {

        javascript: javascriptSnippet,
        python: pythonSnippet,
        go: goSnippet,
        
    };
    return snippets[name] || snippets.javascript;
}
