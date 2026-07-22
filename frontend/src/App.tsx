import { useState } from "react";

import Editor from "./components/Editor";
import Output from "./components/Output";
import ResetButton from "./components/ResetButton";
import RunButton from "./components/RunButton";
import StatusBar from "./components/StatusBar";

import getSnippet from "./snippets/Snippets";

type Language = "javascript" | "python" | "go";

export default function App() {

    const [language, setLanguage] = useState<Language>("go");
    const [code, setCode] = useState(getSnippet("go"));
    const [output, setOutput] = useState("Waiting for execution...");
    const [resetVersion, setResetVersion] = useState(0);

    function handleLanguageChange(nextLanguage: Language) {
        setLanguage(nextLanguage);
        setCode(getSnippet(nextLanguage));
        setOutput("Waiting for execution...");
        setResetVersion((value) => value + 1);
    }

    function handleReset() {
        setCode(getSnippet(language));
        setOutput("Waiting for execution...");
        setResetVersion((value) => value + 1);
    }

    return (
        <main className="container">
            <div className="app-shell">
                <header className="app-header">
                    <div>
                        <h1 className="app-title">LGTM Playground</h1>
                        <p className="app-subtitle">do not trust the terminal</p>
                    </div>
                    <ResetButton onReset={handleReset} />
                    <RunButton code={code} language={language} onResult={setOutput} />
                </header>

                <div className="app-grid">
                    <Editor
                        code={code}
                        language={language}
                        onChange={setCode}
                        onChangeLanguage={handleLanguageChange}
                        resetVersion={resetVersion}
                    />
                    <Output output={output} />
                </div>

                <StatusBar />
            </div>
        </main>
    );
}