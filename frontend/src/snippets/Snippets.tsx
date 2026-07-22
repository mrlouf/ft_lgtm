export default function getSnippet(name: string): string {
    const snippets: Record<string, string> = {
        javascript: `function hello(who = "world") {
    console.log(\`Hello, \${who}!\`)
}`,
        python: `def hello(who="world"):
    print(f"Hello, {who}!")`,
        go: `package main

import "fmt"

func hello(who string) {
    fmt.Printf("Hello, %s!\\n", who)
}`,
    }
    return snippets[name] || snippets.javascript;
}