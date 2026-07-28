import { useState } from "react";

import Editor from "./components/Editor";
import Output from "./components/Output";
import ResetButton from "./components/ResetButton";
import RunButton from "./components/RunButton";
import StatusBar from "./components/StatusBar";

import getSnippet from "./snippets/Snippets";

type Language = "javascript" | "python" | "go";
type Status = "Ready" | "Running" | "Completed" | "Error";

export default function App() {
    const [language, setLanguage] = useState<Language>("go");
    const [code, setCode] = useState(getSnippet("go"));
    const [output, setOutput] = useState("Waiting for execution...");
    const [status, setStatus] = useState<Status>("Ready");
    const [resetVersion, setResetVersion] = useState(0);
    const [cid, setCid] = useState("");

    function handleStatusChange(nextStatus: Status, nextCid: string) {
        setStatus(nextStatus);
        setCid(nextCid);
    }

    function handleLanguageChange(nextLanguage: Language) {
        setLanguage(nextLanguage);
        setCode(getSnippet(nextLanguage));
        setOutput("Waiting for execution...");
        setStatus("Ready");
        setResetVersion((value) => value + 1);
        setCid("");
    }

    function handleEnterCID(nextCid: string) {
        const gatewayUrl = import.meta.env.VITE_IPFS_GATEWAY_URL ?? "http://localhost:8080";

        fetch(`${gatewayUrl}/ipfs/${nextCid}`, {
            method: "GET",
            headers: {
                "Content-Type": "text/plain",
            },
        })
            .then((response) => {
                if (!response.ok) {
                    throw new Error(`Failed to fetch CID ${nextCid}`);
                }
                return response.text();
            })
            .then((source) => {
                setCode(source);
                setCid(nextCid);
                setOutput("Loaded source from IPFS");
                setStatus("Ready");
                setResetVersion((value) => value + 1);
            })
            .catch((error) => {
                console.error("Error fetching snippet:", error);
                setOutput(`Error loading CID: ${error.message}`);
                setStatus("Error");
            });
    }

    function handleReset() {
        setCode(getSnippet(language));
        setOutput("Waiting for execution...");
        setStatus("Ready");
        setResetVersion((value) => value + 1);
        setCid("");
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
                    <RunButton
                        code={code}
                        language={language}
                        onResult={setOutput}
                        onStatusChange={handleStatusChange}
                    />
                </header>

                <Editor
                    code={code}
                    language={language}
                    cid={cid}
                    onChange={setCode}
                    onChangeLanguage={handleLanguageChange}
                    onEnterCID={handleEnterCID}
                    resetVersion={resetVersion}
                />

                <Output output={output} />

                <StatusBar status={status} snippetCID={cid} />
            </div>
        </main>
    );
}