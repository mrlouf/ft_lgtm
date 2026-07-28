type Language = "javascript" | "python" | "go";
type Status = "Ready" | "Running" | "Completed" | "Error";

type RunButtonProps = {
    code: string;
    language: Language;
    onResult: (output: string) => void;
    onStatusChange: (status: Status) => void;
};

export default function RunButton({ code, language, onResult, onStatusChange }: RunButtonProps) {
    function handleRun() {
        onStatusChange("Running");
        onResult("Running code...");

        fetch("/api/run", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                code,
                language,
            }),
        })
            .then((response) => response.json())
            .then((data) => {
                const output = JSON.stringify(data, null, 2);
                console.log("Code execution result:", output);

                if (data.status === "failed") {
                    onStatusChange("Error");
                    onResult(`Error: ${data.error}`);
                } else {

                    data.stdout = data.stdout || "";
                    data.stderr = data.stderr || "";

                    const resultOutput = `Output:\n${data.stdout}\n\n`;
                    if (data.stderr) {
                        const errOutput = `Errors:\n${data.stderr}\n\n`;
                        onResult(resultOutput + errOutput);
                    } else {
                        onResult(resultOutput);
                    }
                }

                onStatusChange("Completed");

                const resultOutput = `${data.stdout}\n\n`;
                if (data.stderr) {
                    const errOutput = `Errors:\n${data.stderr}\n\n`;
                    onResult(resultOutput + errOutput);
                } else {
                    onResult(resultOutput);
                }
            })
            .catch((error) => {
                console.error("Error running code:", error);
                onStatusChange("Error");
                onResult(`Error running code: ${error.message}`);
            });
    }

    return (
        <div className="controls">
            <button className="run-button" onClick={handleRun}>
                Run Code
            </button>
        </div>
    );
}