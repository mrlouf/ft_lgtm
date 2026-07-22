type Language = "javascript" | "python" | "go";

type RunButtonProps = {
    code: string;
    language: Language;
    onResult: (output: string) => void;
};

export default function RunButton({ code, language, onResult }: RunButtonProps) {
    function handleRun() {

        onResult("Running code...");

        fetch("http://localhost:4242/api/run", {
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
                onResult(JSON.stringify(data, null, 2));
            })
            .catch((error) => {
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