type Language = "javascript" | "python" | "go";

type RunButtonProps = {
    code: string;
    language: Language;
};

export default function RunButton({ code, language }: RunButtonProps) {
    function handleRun() {
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
                console.log("Run result:", data);
            })
            .catch((error) => {
                console.error("Error running code:", error);
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