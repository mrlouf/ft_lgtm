export default function RunButton() {
    function handleRun() {
        console.log("Running code...");

        try {

            const res = fetch("http://localhost:4242/api/run", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    code: "console.log('Hello, World!');",
                    language: "javascript",
                }),
            })
                .then((response) => response.json())
                .then((data) => {
                    console.log("Run result:", data);
                })
                .catch((error) => {
                    console.error("Error running code:", error);
                });
        } catch (error) {
            console.error("Error running code:", error);
        }
    }

    return (
        <div className="controls">
            <button className="run-button" onClick={handleRun}>
                Run Code
            </button>
        </div>
    );
}