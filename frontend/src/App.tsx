import { useState } from "react";

import Editor from "./components/Editor";
import Output from "./components/Output";
import ResetButton from "./components/ResetButton";
import RunButton from "./components/RunButton";
import StatusBar from "./components/StatusBar";

import getSnippet from "./snippets/Snippets";

export default function App() {

    const [code, setCode] = useState(getSnippet("javascript"));
    const [resetVersion, setResetVersion] = useState(0);

    function handleReset() {
        setCode(getSnippet("javascript"));
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
                    <RunButton code={code} />
                </header>

                <div className="app-grid">
                    <Editor resetVersion={resetVersion} code={code} onChange={setCode} />
                    <Output />
                </div>

                <StatusBar />
            </div>
        </main>
    );
}