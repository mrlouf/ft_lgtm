import { useState } from "react";

import Editor from "./components/Editor";
import Output from "./components/Output";
import ResetButton from "./components/ResetButton";
import RunButton from "./components/RunButton";
import StatusBar from "./components/StatusBar";

export default function App() {
    const [resetVersion, setResetVersion] = useState(0);

    function handleReset() {
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
                    <RunButton />
                </header>

                <div className="app-grid">
                    <Editor resetVersion={resetVersion} />
                    <Output />
                </div>

                <StatusBar />
            </div>
        </main>
    );
}