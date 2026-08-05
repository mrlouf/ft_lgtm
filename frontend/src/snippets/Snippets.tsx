const goSnippet = `package main

func main() {
    println("Hello, world!")
}`;

const pythonSnippet = `def hello():
    print("Hello, world!")`;

const javascriptSnippet = `function hello() {
    console.log("Hello, world!")
}
    
hello();`;

const cSnippet = `#include <stdio.h>

int main() {
    printf("Hello, world!\\n");
    return (0);
}`;

export default function getSnippet(name: string): string {

    const snippets: Record<string, string> = {

        javascript: javascriptSnippet,
        python: pythonSnippet,
        go: goSnippet,
        c: cSnippet,
        
    };
    return snippets[name] || snippets.javascript;
}
