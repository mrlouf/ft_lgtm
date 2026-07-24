type Language = "javascript" | "python" | "go";

type RunButtonProps = {
    code: string;
    language: Language;
    onResult: (output: string) => void;
};

export default function RunButton({ code, language, onResult }: RunButtonProps) {
    function handleRun() {

        onResult("Running code...");

        const apiUrl = "http://backend:4242";

        fetch("http://lgtm.local/api/run", {
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
                    onResult(`Error: ${data.error}`);
                } else {

                    const resultOutput = `Output:\n${data.stdout}\n\n`;
                    if (data.stderr) {
                        const errOutput = `Errors:\n${data.stderr}\n\n`;
                        onResult(resultOutput + errOutput);
                    } else {
                        onResult(resultOutput);
                    }
                }

            })
            .catch((error) => {

                console.error("Error running code:", error);

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