const goSnippet = `package main

import "fmt"

func main() {
    fmt.Printf("Hello, world!\n")
}`;

const pythonSnippet = `def hello():
    print("Hello, world!")`;

const javascriptSnippet = `function hello() {
    console.log("Hello, world!")
}`;

export default function getSnippet(name: string): string {

    const snippets: Record<string, string> = {

        javascript: javascriptSnippet,
        python: pythonSnippet,
        go: goSnippet,
        
    };
    return snippets[name] || snippets.javascript;
}
