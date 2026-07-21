import { useState } from "react";

import Editor from "./components/Editor";
import Output from "./components/Output";
import ResetButton from "./components/ResetButton";
import RunButton from "./components/RunButton";
import StatusBar from "./components/StatusBar";

export default function App() {
    const [code, setCode] = useState(`// do not trust the terminal
`);
    const [resetVersion, setResetVersion] = useState(0);

    function handleReset() {
        setCode(`// do not trust the terminal
`);
        setResetVersion((value) => value + 1);
    }

    return (
        <main className="container">
            <div className="app-shell">
                <header className="app-header">
                    <div>
                        <h1 className="app-title">LGTM Playground</h1>
                        <p className="app-subtitle">uplink-mode interface</p>
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